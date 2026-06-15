package repository

import (
	"context"
	"time"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
)

type ClienteMoroso struct {
	ID               string    `json:"id"`
	NombreCompleto   string    `json:"nombre_completo"`
	Telefono         string    `json:"telefono"`
	PlanID           string    `json:"plan_id"`
	TerrenoID        string    `json:"terreno_id"`
	TerrenoClave     string    `json:"terreno_clave"`
	PeriodoID        string    `json:"periodo_id"`
	NumeroPeriodo    int       `json:"numero_periodo"`
	MontoEsperado    float64   `json:"monto_esperado"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	DiasRetraso      int       `json:"dias_retraso"`
}

func GetClientesMorosos(parcelaID string) ([]ClienteMoroso, error) {
	query := `
		SELECT 
			c.id, c.nombre_completo, COALESCE(c.telefono, ''), 
			pl.id as plan_id, t.id as terreno_id, t.clave as terreno_clave,
			p.id as periodo_id, p.numero_periodo, p.monto_esperado, p.fecha_vencimiento,
			CURRENT_DATE - p.fecha_vencimiento as dias_retraso
		FROM periodos_pago p
		JOIN planes_pago pl ON p.plan_id = pl.id
		JOIN clientes c ON pl.cliente_id = c.id
		JOIN terrenos t ON pl.terreno_id = t.id
		WHERE pl.parcela_id = $1 
		  AND p.estado = 'pendiente' 
		  AND p.fecha_vencimiento < CURRENT_DATE
		  AND t.estado != 'disponible'
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
			&m.PlanID, &m.TerrenoID, &m.TerrenoClave, &m.PeriodoID, 
			&m.NumeroPeriodo, &m.MontoEsperado, &m.FechaVencimiento, &m.DiasRetraso,
		); err != nil {
			return nil, err
		}
		morosos = append(morosos, m)
	}
	return morosos, nil
}
