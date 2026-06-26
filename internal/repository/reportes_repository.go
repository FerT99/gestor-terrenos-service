package repository

import (
	"context"
	"time"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
)

type ClienteMoroso struct {
	ID               string     `json:"id"`
	NombreCompleto   string     `json:"nombre_completo"`
	Telefono         string     `json:"telefono"`
	PlanID           string     `json:"plan_id"`
	TerrenoID        string     `json:"terreno_id"`
	TerrenoClave     string     `json:"terreno_clave"`
	DiasRetraso      int        `json:"dias_retraso"`
	UltimoAbonoFecha *time.Time `json:"ultimo_abono_fecha"`
}

func GetClientesMorosos(parcelaID string) ([]ClienteMoroso, error) {
	query := `
		SELECT 
			c.id, c.nombre_completo, COALESCE(c.telefono, ''), 
			pl.id as plan_id, t.id as terreno_id, t.clave as terreno_clave,
			MAX(CURRENT_DATE - p.fecha_vencimiento) as dias_retraso,
			(
			  SELECT MAX(fecha_pago)
			  FROM abonos a
			  JOIN periodos_pago pp ON a.periodo_pago_id = pp.id
			  WHERE pp.plan_id = pl.id AND a.fecha_pago <= CURRENT_DATE
			) as ultimo_abono_fecha
		FROM periodos_pago p
		JOIN planes_pago pl ON p.plan_id = pl.id
		JOIN clientes c ON pl.cliente_id = c.id
		JOIN terrenos t ON pl.terreno_id = t.id
		WHERE pl.parcela_id = $1 
		  AND p.estado = 'pendiente' 
		  AND p.fecha_vencimiento < CURRENT_DATE
		  AND t.estado != 'disponible'
		  AND (
			  SELECT COALESCE(SUM(monto_pagado), 0) 
			  FROM abonos a 
			  JOIN periodos_pago pp ON a.periodo_pago_id = pp.id 
			  WHERE pp.plan_id = pl.id
		  ) < t.precio_lista
		  AND NOT EXISTS (
			  SELECT 1 
			  FROM abonos a
			  JOIN periodos_pago pp ON a.periodo_pago_id = pp.id
			  WHERE pp.plan_id = pl.id 
			    AND EXTRACT(MONTH FROM a.fecha_pago) = EXTRACT(MONTH FROM CURRENT_DATE)
			    AND EXTRACT(YEAR FROM a.fecha_pago) = EXTRACT(YEAR FROM CURRENT_DATE)
		  )
		  AND (
			  SELECT COALESCE(SUM(monto_esperado), 0) FROM periodos_pago WHERE plan_id = pl.id
		  ) > (
			  SELECT COALESCE(SUM(monto_pagado), 0) 
			  FROM abonos a 
			  JOIN periodos_pago pp ON a.periodo_pago_id = pp.id 
			  WHERE pp.plan_id = pl.id
		  )
		GROUP BY c.id, c.nombre_completo, c.telefono, pl.id, t.id, t.clave
		ORDER BY dias_retraso DESC
	`
	
	rows, err := database.DB.Query(context.Background(), query, parcelaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var morosos []ClienteMoroso
	for rows.Next() {
		var m ClienteMoroso
		if err := rows.Scan(
			&m.ID, &m.NombreCompleto, &m.Telefono, 
			&m.PlanID, &m.TerrenoID, &m.TerrenoClave,
			&m.DiasRetraso, &m.UltimoAbonoFecha,
		); err != nil {
			return nil, err
		}
		morosos = append(morosos, m)
	}
	return morosos, nil
}
