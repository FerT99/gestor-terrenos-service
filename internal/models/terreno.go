package models

import "time"

// Terreno representa un lote/terreno en la base de datos
type Terreno struct {
	ID           string    `json:"id"`
	Clave        string    `json:"clave"`
	Nombre       *string   `json:"nombre"`
	Fase         *string   `json:"fase"`
	SuperficieM2 float64   `json:"superficie_m2"`
	PrecioLista  float64   `json:"precio_lista"`
	Propietario  *string   `json:"propietario"`
	Estado       string    `json:"estado"`
	Coordenadas  *string   `json:"coordenadas"`
	Notas        *string   `json:"notas"`
	CreatedAt    time.Time `json:"created_at"`
}

// TerrenoInput es el payload para crear o actualizar un terreno
type TerrenoInput struct {
	Clave        string  `json:"clave"`
	Nombre       string  `json:"nombre"`
	Fase         string  `json:"fase"`
	SuperficieM2 float64 `json:"superficie_m2"`
	PrecioLista  float64 `json:"precio_lista"`
	Propietario  string  `json:"propietario"`
	Estado       string  `json:"estado"`
	Coordenadas  string  `json:"coordenadas"`
	Notas        string  `json:"notas"`
}
