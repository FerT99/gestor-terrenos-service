package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func GetUsuarios(c *fiber.Ctx) error {
	usuarios, err := repository.GetAllUsuarios()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener usuarios: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": usuarios})
}

func CreateOrUpdateUsuario(c *fiber.Ctx) error {
	var input models.UsuarioInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}

	if input.ID == "" || input.Email == "" || input.NombreCompleto == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID, Email y Nombre son requeridos",
		})
	}

	if input.Rol == "" {
		input.Rol = "vendedor"
	}

	usuario, err := repository.CreateOrUpdateUsuario(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al crear/actualizar usuario: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": usuario})
}

func GetMe(c *fiber.Ctx) error {
	userID := c.Get("X-User-Id")
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-User-Id requerido",
		})
	}
	usuario, err := repository.GetUsuarioByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Usuario no encontrado",
		})
	}
	return c.JSON(fiber.Map{"error": false, "data": usuario})
}
