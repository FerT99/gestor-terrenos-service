package repository

import (
	"context"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

func CreatePlanPago(input models.PlanPagoInput) (models.PlanPago, error) {
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return models.PlanPago{}, err
	}
	defer tx.Rollback(context.Background())

	queryPlan := `
		INSERT INTO planes_pago (parcela_id, terreno_id, cliente_id, monto_total, enganche, plazos, tasa_interes, fecha_inicio, estado)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, parcela_id, terreno_id, cliente_id, monto_total, enganche, plazos, tasa_interes, fecha_inicio, estado, created_at
	`
	var plan models.PlanPago
	err = tx.QueryRow(
		context.Background(),
		queryPlan,
		input.ParcelaID,
		input.TerrenoID,
		input.ClienteID,
		input.MontoTotal,
		input.Enganche,
		input.Plazos,
		input.TasaInteres,
		input.FechaInicio,
		"Activo",
	).Scan(
		&plan.ID, &plan.ParcelaID, &plan.TerrenoID, &plan.ClienteID,
		&plan.MontoTotal, &plan.Enganche, &plan.Plazos,
		&plan.TasaInteres, &plan.FechaInicio, &plan.Estado, &plan.CreatedAt,
	)
	if err != nil {
		return models.PlanPago{}, err
	}

	// Obtener el nombre del cliente
	var nombreCliente string
	err = tx.QueryRow(context.Background(), "SELECT nombre_completo FROM clientes WHERE id = $1", input.ClienteID).Scan(&nombreCliente)
	if err != nil {
		return models.PlanPago{}, err
	}

	// Update terreno status to "vendido" and set propietario
	queryUpdateTerreno := `UPDATE terrenos SET estado = 'vendido', propietario = $2 WHERE id = $1`
	_, err = tx.Exec(context.Background(), queryUpdateTerreno, input.TerrenoID, nombreCliente)
	if err != nil {
		return models.PlanPago{}, err
	}

	// Generar periodos de pago
	montoFinanciar := input.MontoTotal - input.Enganche
	interesTotal := montoFinanciar * (input.TasaInteres / 100)
	totalAPagar := montoFinanciar + interesTotal
	mensualidad := 0.0
	if input.Plazos > 0 {
		mensualidad = totalAPagar / float64(input.Plazos)
	}

	queryPeriodo := `
		INSERT INTO periodos_pago (plan_id, numero_periodo, monto_esperado, fecha_vencimiento, estado, mora_aplicada)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for i := 1; i <= input.Plazos; i++ {
		fechaVenc := input.FechaInicio.AddDate(0, i, 0)
		_, err = tx.Exec(
			context.Background(),
			queryPeriodo,
			plan.ID,
			i,
			mensualidad,
			fechaVenc,
			"pendiente",
			0.0,
		)
		if err != nil {
			return models.PlanPago{}, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return models.PlanPago{}, err
	}

	return plan, nil
}

func GetAllPlanesPago(parcelaID string) ([]models.PlanPago, error) {
	query := `
		SELECT p.id, p.parcela_id, p.terreno_id, p.cliente_id, p.monto_total, p.enganche, p.plazos, p.tasa_interes, p.fecha_inicio, p.estado, p.created_at,
		       c.nombre_completo as cliente_nombre,
		       t.clave as terreno_nombre
		FROM planes_pago p
		LEFT JOIN clientes c ON p.cliente_id = c.id
		LEFT JOIN terrenos t ON p.terreno_id = t.id
		WHERE p.parcela_id = $1
		ORDER BY p.created_at DESC
	`
	rows, err := database.DB.Query(context.Background(), query, parcelaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var planes []models.PlanPago
	for rows.Next() {
		var p models.PlanPago
		if err := rows.Scan(
			&p.ID, &p.ParcelaID, &p.TerrenoID, &p.ClienteID, &p.MontoTotal, &p.Enganche, &p.Plazos, &p.TasaInteres, &p.FechaInicio, &p.Estado, &p.CreatedAt,
			&p.ClienteNombre, &p.TerrenoNombre,
		); err != nil {
			return nil, err
		}
		planes = append(planes, p)
	}
	return planes, nil
}

func GetPeriodosByPlan(planID string) ([]models.PeriodoPago, error) {
	query := `
		SELECT id, plan_id, numero_periodo, monto_esperado, fecha_vencimiento, estado, mora_aplicada, created_at 
		FROM periodos_pago 
		WHERE plan_id = $1 
		ORDER BY numero_periodo ASC
	`
	rows, err := database.DB.Query(context.Background(), query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periodos []models.PeriodoPago
	for rows.Next() {
		var p models.PeriodoPago
		if err := rows.Scan(&p.ID, &p.PlanID, &p.NumeroPeriodo, &p.MontoEsperado, &p.FechaVencimiento, &p.Estado, &p.MoraAplicada, &p.CreatedAt); err != nil {
			return nil, err
		}
		periodos = append(periodos, p)
	}
	return periodos, nil
}
