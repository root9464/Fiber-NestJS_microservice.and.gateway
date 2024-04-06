package routes_auth

import (
	"context"
	"log"
	"sync"
	"time"

	pb "root/proto/out"

	"github.com/gofiber/fiber/v2"
)

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

func Login(microserviceClient pb.UserServiceClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		message := pb.LoginReq{}
		if err := c.BodyParser(&message); err != nil {
			log.Fatal("ошибка: не удалось прочитать данные из body", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		ch := make(chan error, 1)

		go func() {
			mu.Lock()
			defer mu.Unlock()
			_, err := microserviceClient.Login(ctx, &message)
			ch <- err
		}()

		select {
		case err := <-ch:
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		case <-ctx.Done():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "request timeout"})
		}

		return c.JSON(fiber.Map{"status": "ok"})
	}
}

func Register(microserviceClient pb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var message pb.RegisterReq
		if err := c.BodyParser(&message); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		ch := make(chan error)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		go func() {
			mu.Lock()
			defer mu.Unlock()
			defer cancel()
			_, err := microserviceClient.Register(ctx, &message)
			ch <- err
		}()

		select {
		case err := <-ch:
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		case <-ctx.Done():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "request timeout"})
		}

		return c.JSON(fiber.Map{"message": "Регистрация пользователя запущена"})
	}
}

func LogoutService(microserviceClient pb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get all cookies
		cookies := c.Cookies("token")
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    "",
			HTTPOnly: true,
			SameSite: "",
			Expires:  time.Now().Add(time.Duration(0) * time.Second),
			Path:     "/",
		})
		c.ClearCookie(cookies)
		dto := pb.LogoutReq{
			Agent: c.Get("User-Agent"),
		}

		ctx := context.Background()

		ch := make(chan error, 1)

		go func() {
			mu.Lock()
			defer mu.Unlock()

			_, err := microserviceClient.Logout(ctx, &dto)
			ch <- err
		}()

		select {
		case err := <-ch:
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err.Error()})
			}
		default:
			return c.JSON(fiber.Map{"status": "ok"})
		}

		return c.JSON(fiber.Map{"status": "ok"})
	}
}

func VerifyCode(microserviceClient pb.UserServiceClient) fiber.Handler {

	return func(c *fiber.Ctx) error {
		message := pb.VerifyCodeReq{
			Agent: c.Get("User-Agent"),
			Body:  &pb.VerifyCodeBody{},
		}
		if err := c.BodyParser(&message); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		var response *pb.AccessTokenRes
		var err error

		wg.Add(1)
		go func() {
			defer wg.Done()
			response, err = microserviceClient.VerifyCode(ctx, &message)
		}()

		wg.Wait()
		cancel()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    response.AccessToken.Token,
			HTTPOnly: true,
			SameSite: "lax",
			Expires:  time.Now().Add(time.Duration(response.AccessToken.Exp) * time.Second),
			Path:     "/",
		})

		return c.JSON(fiber.Map{"message": "Код подтверждения отправлен"})
	}
}
