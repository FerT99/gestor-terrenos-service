package models

import "time"

type Cliente struct {
	ID             string    `json:"id"`
	ParcelaID      string    `json:"parcela_id"`
	NombreCompleto string    `json:"nombre_completo"`
	Email          *string   `json:"email"`
	Telefono       *string   `json:"telefono"`
	Direccion      *string   `json:"direccion"`
	Estado         string    `json:"estado"`
	CreatedAt      time.Time `json:"created_at"`
}

type ClienteInput struct {
	ParcelaID      string `json:"parcela_id"`
	NombreCompleto string `json:"nombre_completo" validate:"required"`
	Email          string `json:"email"`
	Telefono       string `json:"telefono"`
	Direccion      string `json:"direccion"`
	Estado         string `json:"estado"` // ej: Activo, Pendiente
}
