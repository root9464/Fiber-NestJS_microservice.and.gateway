package main

import (
	"log"
	routes_auth "root/routes/auth"
	"sync"

	pb "root/proto/out"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	once               sync.Once
	microserviceClient pb.UserServiceClient
	wg                 sync.WaitGroup
)

func initGRPCClient() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := grpc.NewClient("127.0.0.1:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("не получилось соединиться: %v", err)
		}
		microserviceClient = pb.NewUserServiceClient(conn)
	}()

	wg.Wait()
}

func AllRoutes(app *fiber.App, microserviceClient pb.UserServiceClient) {
	auth := app.Group("/auth")
	auth.Post("/login", routes_auth.Login(microserviceClient))
	auth.Post("/register", routes_auth.Register(microserviceClient))
	auth.Post("/code", routes_auth.VerifyCode(microserviceClient))
	auth.Get("/logout/:id", routes_auth.LogoutService(microserviceClient))
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	once.Do(initGRPCClient)

	AllRoutes(app, microserviceClient)

	log.Fatal(app.Listen(":6069"))
}
