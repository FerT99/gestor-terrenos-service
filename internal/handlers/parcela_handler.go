package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func GetParcelas(c *fiber.Ctx) error {
	parcelas, err := repository.GetAllParcelas()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener parcelas: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": parcelas})
}

func GetParcelaByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID requerido",
		})
	}
	parcela, err := repository.GetParcelaByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Parcela no encontrada: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": parcela})
}

func CreateParcela(c *fiber.Ctx) error {
	var input models.ParcelaInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}
	if input.Nombre == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "El nombre de la parcela es requerido",
		})
	}

	parcela, err := repository.CreateParcela(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al crear parcela: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "data": parcela})
}

func UpdateParcela(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID requerido",
		})
	}

	var input models.ParcelaInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}

	parcela, err := repository.UpdateParcela(id, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al actualizar parcela: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": parcela})
}

func DeleteParcela(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID requerido",
		})
	}
	if err := repository.DeleteParcela(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al eliminar parcela: " + err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
