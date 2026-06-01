package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type RegisterVendedorInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nombre   string `json:"nombre"`
	Rol      string `json:"rol"`
}

func RegisterVendedor(c *fiber.Ctx) error {
	// Verificar que sea admin
	rol := c.Get("X-User-Role")
	if rol != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Solo los administradores pueden crear vendedores",
		})
	}

	var input RegisterVendedorInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos",
		})
	}

	if input.Email == "" || input.Password == "" || input.Nombre == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Email, nombre y contraseña son requeridos",
		})
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" || serviceRoleKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Configuración de Supabase no encontrada en el servidor",
		})
	}

	// Crear el body para la API de Supabase Admin
	requestBody, _ := json.Marshal(map[string]interface{}{
		"email":         input.Email,
		"password":      input.Password,
		"email_confirm": true, // Para que puedan iniciar sesión de inmediato
	})

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/admin/users", supabaseURL), bytes.NewBuffer(requestBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Error creando petición HTTP"})
	}

	req.Header.Set("apikey", serviceRoleKey)
	req.Header.Set("Authorization", "Bearer "+serviceRoleKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Error al contactar con Supabase"})
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error":   true,
			"message": "Error de Supabase al crear usuario: " + string(bodyBytes),
		})
	}

	// Parsear la respuesta de Supabase para obtener el ID
	var supaResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(bodyBytes, &supaResp); err != nil || supaResp.ID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "No se pudo obtener el ID del usuario desde Supabase",
		})
	}

	rolAInsertar := "vendedor"
	if input.Rol == "admin" {
		rolAInsertar = "admin"
	}

	// Insertar en nuestra tabla de usuarios
	usuarioInput := models.UsuarioInput{
		ID:             supaResp.ID,
		NombreCompleto: input.Nombre,
		Email:          input.Email,
		Rol:            rolAInsertar,
	}

	usuario, err := repository.CreateOrUpdateUsuario(usuarioInput)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Usuario creado en Auth, pero falló al guardar en DB: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "data": usuario})
}

// ... rest of the handlers ...
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
