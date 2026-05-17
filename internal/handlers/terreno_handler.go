package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func GetTerrenos(c *fiber.Ctx) error {
	terrenos, err := repository.GetAllTerrenos()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener terrenos: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  terrenos,
	})
}
