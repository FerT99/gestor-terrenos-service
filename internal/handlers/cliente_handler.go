package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func GetClientes(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}
	clientes, err := repository.GetAllClientes(parcelaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener clientes: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": clientes})
}

func GetClienteByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID requerido",
		})
	}
	cliente, err := repository.GetClienteByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Cliente no encontrado: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": cliente})
}

func CreateCliente(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}

	var input models.ClienteInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}
	input.ParcelaID = parcelaID

	if input.NombreCompleto == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "El nombre completo es requerido",
		})
	}

	cliente, err := repository.CreateCliente(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al crear cliente: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "data": cliente})
}

func UpdateCliente(c *fiber.Ctx) error {
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

	var input models.ClienteInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}
	input.ParcelaID = parcelaID

	if input.NombreCompleto == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "El nombre completo es requerido",
		})
	}

	cliente, err := repository.UpdateCliente(id, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al actualizar cliente: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "data": cliente})
}

func DeleteCliente(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID requerido",
		})
	}

	if err := repository.DeleteCliente(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al eliminar cliente: " + err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
