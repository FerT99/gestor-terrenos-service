package models

import (
	"time"
)

type Egreso struct {
	ID          string    `json:"id"`
	ParcelaID   string    `json:"parcela_id"`
	Fecha       time.Time `json:"fecha"`
	Concepto    string    `json:"concepto"`
	Monto       float64   `json:"monto"`
	Categoria   string    `json:"categoria"`
	Descripcion string    `json:"descripcion"`
}

type EgresoInput struct {
	Fecha       string  `json:"fecha"`
	Concepto    string  `json:"concepto"`
	Monto       float64 `json:"monto"`
	Categoria   string  `json:"categoria"`
	Descripcion string  `json:"descripcion"`
}
