package models

import "time"

// Parcela representa una entidad que agrupa terrenos, clientes y ventas.
type Parcela struct {
	ID        string    `json:"id"`
	Nombre    string    `json:"nombre"`
	CreatedAt time.Time `json:"created_at"`
}

// ParcelaInput es el payload para crear o actualizar una parcela.
type ParcelaInput struct {
	Nombre string `json:"nombre" validate:"required"`
}
