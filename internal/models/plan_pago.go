package models

import (
	"time"
)

type PlanPago struct {
	ID           string         `json:"id"`
	ParcelaID    string         `json:"parcela_id"`
	TerrenoID    string         `json:"terreno_id"`
	ClienteID    string         `json:"cliente_id"`
	MontoTotal   float64        `json:"monto_total"`
	Enganche     float64        `json:"enganche"`
	Moneda       string         `json:"moneda"`
	Plazos       int            `json:"plazos"`
	TasaInteres  float64        `json:"tasa_interes"`
	FechaInicio  time.Time      `json:"fecha_inicio"`
	Estado       string         `json:"estado"`
	CreatedAt    time.Time      `json:"created_at"`
	
	// Utilizado para joins
	ClienteNombre *string `json:"cliente_nombre,omitempty"`
	TerrenoNombre *string `json:"terreno_nombre,omitempty"`
}

type PeriodoPago struct {
	ID               string    `json:"id"`
	PlanID           string    `json:"plan_id"`
	NumeroPeriodo    int       `json:"numero_periodo"`
	MontoEsperado    float64   `json:"monto_esperado"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Estado           string    `json:"estado"` // pagado, pendiente
	MoraAplicada     float64   `json:"mora_aplicada"`
	CreatedAt        time.Time `json:"created_at"`
}

type PlanPagoInput struct {
	ParcelaID    string    `json:"parcela_id"`
	TerrenoID    string    `json:"terreno_id"`
	ClienteID    string    `json:"cliente_id"`
	MontoTotal   float64   `json:"monto_total"`
	Enganche     float64   `json:"enganche"`
	Moneda       string    `json:"moneda"`
	Plazos       int       `json:"plazos"`
	TasaInteres  float64   `json:"tasa_interes"`
	FechaInicio  time.Time `json:"fecha_inicio"`
}
