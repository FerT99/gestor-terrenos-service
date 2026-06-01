package models

import (
	"time"
)

type Abono struct {
	ID              string         `json:"id"`
	ParcelaID       string         `json:"parcela_id"`
	PeriodoPagoID   string         `json:"periodo_pago_id"`
	NumeroAbono     int            `json:"numero_abono"`
	MontoPagado     float64        `json:"monto_pagado"`
	Moneda          string         `json:"moneda"`
	FechaPago       time.Time      `json:"fecha_pago"`
	MetodoPago      *string        `json:"metodo_pago"`
	ComprobanteURL  *string        `json:"comprobante_url"`
	Notas           *string        `json:"notas"`
	CreatedAt       time.Time      `json:"created_at"`
	TerrenoClave    string         `json:"terreno_clave,omitempty"`
	TerrenoNombre   string         `json:"terreno_nombre,omitempty"`
	ClienteNombre   string         `json:"cliente_nombre,omitempty"`
}

type AbonoInput struct {
	PeriodoPagoID string  `json:"periodo_pago_id"`
	MontoPagado   float64 `json:"monto_pagado"`
	Moneda        string  `json:"moneda"`
	FechaPago     string  `json:"fecha_pago"`
	MetodoPago    string  `json:"metodo_pago"`
	Notas          string  `json:"notas"`
	PerdonarMora   bool    `json:"perdonar_mora"`
	ComprobanteURL string  `json:"comprobante_url"`
}
