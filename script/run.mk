.DEFAULT_GOAL := run

# Определение цели run, которая принимает аргумент path и выполняет go run main.go по указанному пути
run:
    go run $(path)