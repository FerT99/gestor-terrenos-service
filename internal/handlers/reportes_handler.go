package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func GetClientesMorosos(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}

	morosos, err := repository.GetClientesMorosos(parcelaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener morosos: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "data": morosos})
}
