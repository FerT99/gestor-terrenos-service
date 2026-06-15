package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func GetEgresosByParcela(c *fiber.Ctx) error {
	parcelaID := c.Params("parcela_id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID de parcela requerido",
		})
	}

	egresos, err := repository.GetEgresosByParcela(parcelaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener egresos: " + err.Error(),
		})
	}

	if egresos == nil {
		egresos = []models.Egreso{}
	}

	return c.JSON(fiber.Map{"error": false, "data": egresos})
}

func CreateEgreso(c *fiber.Ctx) error {
	parcelaID := c.Params("parcela_id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID de parcela requerido",
		})
	}

	var input models.EgresoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}

	egreso, err := repository.CreateEgreso(parcelaID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al crear egreso: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "data": egreso})
}

func UpdateEgreso(c *fiber.Ctx) error {
	parcelaID := c.Params("parcela_id")
	id := c.Params("id")
	if parcelaID == "" || id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID de parcela y egreso requeridos",
		})
	}

	var input models.EgresoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}

	egreso, err := repository.UpdateEgreso(id, parcelaID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al actualizar egreso: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "data": egreso})
}

func DeleteEgreso(c *fiber.Ctx) error {
	parcelaID := c.Params("parcela_id")
	id := c.Params("id")
	if parcelaID == "" || id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID de parcela y egreso requeridos",
		})
	}

	err := repository.DeleteEgreso(id, parcelaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al eliminar egreso: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "message": "Egreso eliminado correctamente"})
}
