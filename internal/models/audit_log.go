package models

import (
	"time"
)

type AuditLog struct {
	ID            string                 `json:"id"`
	UsuarioNombre string                 `json:"usuario_nombre"`
	Accion        string                 `json:"accion"`
	EntidadTipo   string                 `json:"entidad_tipo"`
	EntidadID     string                 `json:"entidad_id"`
	Detalles      map[string]interface{} `json:"detalles"`
	CreatedAt     time.Time              `json:"created_at"`
}

type AuditLogInput struct {
	UsuarioNombre string                 `json:"usuario_nombre"`
	Accion        string                 `json:"accion"`
	EntidadTipo   string                 `json:"entidad_tipo"`
	EntidadID     string                 `json:"entidad_id"`
	Detalles      map[string]interface{} `json:"detalles"`
}
