package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/models"
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
	return c.JSON(fiber.Map{"error": false, "data": terrenos})
}

func CreateTerreno(c *fiber.Ctx) error {
	var input models.TerrenoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}
	if input.Clave == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "La clave del terreno es requerida",
		})
	}
	if input.Estado == "" {
		input.Estado = "Disponible"
	}

	terreno, err := repository.CreateTerreno(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al crear terreno: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "data": terreno})
}

func UpdateTerreno(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID requerido",
		})
	}

	var input models.TerrenoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}

	terreno, err := repository.UpdateTerreno(id, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al actualizar terreno: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": terreno})
}

func DeleteTerreno(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID requerido",
		})
	}
	if err := repository.DeleteTerreno(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al eliminar terreno: " + err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
