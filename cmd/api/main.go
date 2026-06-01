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
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Parcela-Id, X-User-Id, X-User-Role",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Grupo de rutas protegidas con JWT
	api := app.Group("/api/v1", middleware.Auth())

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "message": "Backend Sahara Lands funcionando correctamente"})
	})

	// Usuarios / Auth
	api.Get("/usuarios", handlers.GetUsuarios)
	api.Post("/usuarios", handlers.CreateOrUpdateUsuario)
	api.Post("/usuarios/vendedores", handlers.RegisterVendedor)
	api.Get("/usuarios/me", handlers.GetMe)

	// CRUD parcelas
	api.Get("/parcelas", handlers.GetParcelas)
	api.Get("/parcelas/:id", handlers.GetParcelaByID)
	api.Post("/parcelas", handlers.CreateParcela)
	api.Put("/parcelas/:id", handlers.UpdateParcela)
	api.Delete("/parcelas/:id", handlers.DeleteParcela)

	// CRUD terrenos
	api.Get("/terrenos", handlers.GetTerrenos)
	api.Get("/terrenos/:id", handlers.GetTerrenoByID)
	api.Post("/terrenos", handlers.CreateTerreno)
	api.Put("/terrenos/:id", handlers.UpdateTerreno)
	api.Delete("/terrenos/:id", handlers.DeleteTerreno)

	// CRUD clientes
	api.Get("/clientes", handlers.GetClientes)
	api.Get("/clientes/:id", handlers.GetClienteByID)
	api.Post("/clientes", handlers.CreateCliente)
	api.Put("/clientes/:id", handlers.UpdateCliente)
	api.Delete("/clientes/:id", handlers.DeleteCliente)

	// Planes de Pago
	api.Post("/planes-pago", handlers.CreatePlanPago)
	api.Get("/planes-pago", handlers.GetPlanesPago)
	api.Get("/planes-pago/:id/periodos", handlers.GetPeriodosPlan)

	// Abonos
	api.Post("/abonos", handlers.CreateAbono)
	api.Get("/abonos", handlers.GetAllAbonos)
	api.Get("/periodos/:periodo_id/abonos", handlers.GetAbonos)

	// Reportes
	api.Get("/reportes/morosos", handlers.GetClientesMorosos)

	// Historial de Actividad
	api.Get("/audit-logs", handlers.GetAuditLogs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Iniciando servidor en puerto %s...\n", port)
	log.Fatal(app.Listen(":" + port))
}
