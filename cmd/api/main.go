package main

import (
	"log"
	"os"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró archivo .env, usando variables del sistema")
	}

	// Conectar a la base de datos
	database.ConnectDB()

	// Inicializar Fiber
	app := fiber.New()

	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // En producción, cambiar por el dominio del frontend
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Grupo de rutas para la API v1
	api := app.Group("/api/v1")

	// Endpoints
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "message": "Backend Sahara Lands funcionando correctamente"})
	})

	api.Get("/terrenos", handlers.GetTerrenos)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Iniciando servidor en puerto %s...\n", port)
	log.Fatal(app.Listen(":" + port))
}
