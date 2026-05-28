package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func GetTerrenos(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}
	vendedorID := c.Get("X-User-Id")
	role := c.Get("X-User-Role")
	var vID *string
	if role == "vendedor" && vendedorID != "" {
		vID = &vendedorID
	}

	terrenos, err := repository.GetAllTerrenos(parcelaID, vID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener terrenos: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": terrenos})
}

func GetTerrenoByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID requerido",
		})
	}
	terreno, err := repository.GetTerrenoByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Terreno no encontrado: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": terreno})
}

func CreateTerreno(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}

	var input models.TerrenoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}
	input.ParcelaID = parcelaID
	if input.Clave == "" {
		nextClave, err := repository.GetNextClave(parcelaID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Error al generar clave automática: " + err.Error(),
			})
		}
		input.Clave = nextClave
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
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}

	var input models.TerrenoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}
	input.ParcelaID = parcelaID

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
