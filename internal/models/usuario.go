package models

import "time"

type Usuario struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	NombreCompleto string    `json:"nombre_completo"`
	Rol            string    `json:"rol"`
	CreatedAt      time.Time `json:"created_at"`
}

type UsuarioInput struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	NombreCompleto string `json:"nombre_completo"`
	Rol            string `json:"rol"`
}
