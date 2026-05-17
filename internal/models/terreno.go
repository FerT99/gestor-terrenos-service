package models

import (
	"time"
)

type Terreno struct {
	ID                 string     `json:"id"`
	Clave              string     `json:"clave"`
	Nombre             string     `json:"nombre"`
	SuperficieM2       float64    `json:"superficie_m2"`
	Precio             float64    `json:"precio"`
	PropietarioFamiliar string    `json:"propietario_familiar"`
	Estado             string     `json:"estado"`
	Notas              string     `json:"notas"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
