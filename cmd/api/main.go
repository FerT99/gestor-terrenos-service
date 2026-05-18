package main

import (
	"log"
	"os"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/handlers"
	"github.com/FerT99/gestor-terrenos-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró archivo .env, usando variables del sistema")
	}

	database.ConnectDB()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Grupo de rutas protegidas con JWT
	api := app.Group("/api/v1", middleware.Auth())

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "message": "Backend Sahara Lands funcionando correctamente"})
	})

	// CRUD terrenos
	api.Get("/terrenos", handlers.GetTerrenos)
	api.Post("/terrenos", handlers.CreateTerreno)
	api.Put("/terrenos/:id", handlers.UpdateTerreno)
	api.Delete("/terrenos/:id", handlers.DeleteTerreno)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Iniciando servidor en puerto %s...\n", port)
	log.Fatal(app.Listen(":" + port))
}
