package handlers

import (
	"context"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetAuditLogs(c *fiber.Ctx) error {
	query := `
		SELECT id, usuario_nombre, accion, entidad_tipo, entidad_id, detalles, created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener historial: " + err.Error(),
		})
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.UsuarioNombre, &l.Accion, &l.EntidadTipo, &l.EntidadID, &l.Detalles, &l.CreatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   true,
				"message": "Error al leer historial: " + err.Error(),
			})
		}
		logs = append(logs, l)
	}

	return c.JSON(fiber.Map{"error": false, "data": logs})
}
