package main

import (
    "go-serve-pos/database"
    "go-serve-pos/routes"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    database.ConnectDB()

    app := fiber.New()

    app.Use(cors.New())
    app.Use(logger.New())
    app.Static("/images", "./uploads")

    routes.SetupRoutes(app)

    app.Listen(":3000")
}