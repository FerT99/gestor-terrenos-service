package repository

import (
	"context"
	"log"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

// LogAction guarda en la base de datos una bitácora de actividad.
// Al ser un proceso de auditoría, se recomienda correr en una goroutine o no fallar la petición si esto falla (best effort).
func LogAction(input models.AuditLogInput) {
	query := `
		INSERT INTO audit_logs (usuario_nombre, accion, entidad_tipo, entidad_id, detalles)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := database.DB.Exec(
		context.Background(),
		query,
		input.UsuarioNombre,
		input.Accion,
		input.EntidadTipo,
		input.EntidadID,
		input.Detalles,
	)

	if err != nil {
		log.Printf("ERROR: No se pudo guardar el audit log para accion %s en %s: %v\n", input.Accion, input.EntidadTipo, err)
	}
}
