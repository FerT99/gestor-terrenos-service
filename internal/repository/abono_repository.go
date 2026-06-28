package repository

import (
	"context"
	"errors"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

func CreateAbono(parcelaID string, input models.AbonoInput) (models.Abono, error) {
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return models.Abono{}, err
	}
	defer tx.Rollback(context.Background())

	// 1. Obtener el periodo de pago
	var periodo models.PeriodoPago
	queryPeriodo := `SELECT id, plan_id, numero_periodo, monto_esperado, fecha_vencimiento, estado, mora_aplicada FROM periodos_pago WHERE id = $1`
	err = tx.QueryRow(context.Background(), queryPeriodo, input.PeriodoPagoID).Scan(
		&periodo.ID, &periodo.PlanID, &periodo.NumeroPeriodo, &periodo.MontoEsperado, 
		&periodo.FechaVencimiento, &periodo.Estado, &periodo.MoraAplicada,
	)
	if err != nil {
		return models.Abono{}, errors.New("periodo de pago no encontrado")
	}

	if periodo.Estado == "pagado" {
		return models.Abono{}, errors.New("este periodo ya está pagado")
	}

	// 2. Mora manual provista desde el frontend
	moraAplicada := input.MoraAplicada

	// 2.5 Obtener el siguiente numero_abono
	var maxAbono *int
	queryMax := `
		SELECT MAX(numero_abono) 
		FROM abonos a
		JOIN periodos_pago pp ON a.periodo_pago_id = pp.id
		WHERE pp.plan_id = $1
	`
	err = tx.QueryRow(context.Background(), queryMax, periodo.PlanID).Scan(&maxAbono)
	numeroAbono := 1
	if maxAbono != nil {
		numeroAbono = *maxAbono + 1
	}

	// 3. Registrar el Abono
	queryAbono := `
		INSERT INTO abonos (parcela_id, periodo_pago_id, numero_abono, monto_pagado, moneda, tipo_cambio, fecha_pago, metodo_pago, comprobante_url, notas)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, parcela_id, periodo_pago_id, numero_abono, monto_pagado, moneda, tipo_cambio, fecha_pago, metodo_pago, comprobante_url, notas, created_at
	`
	var abono models.Abono
	var compURL *string
	if input.ComprobanteURL != "" {
		compURL = &input.ComprobanteURL
	}
	var tipoCambio *float64
	if input.TipoCambio > 0 {
		tipoCambio = &input.TipoCambio
	}
	err = tx.QueryRow(
		context.Background(),
		queryAbono,
		parcelaID,
		input.PeriodoPagoID,
		numeroAbono,
		input.MontoPagado,
		input.Moneda,
		tipoCambio,
		input.FechaPago,
		input.MetodoPago,
		compURL,
		input.Notas,
	).Scan(
		&abono.ID, &abono.ParcelaID, &abono.PeriodoPagoID, &abono.NumeroAbono, &abono.MontoPagado, 
		&abono.Moneda, &abono.TipoCambio, &abono.FechaPago, &abono.MetodoPago, &abono.ComprobanteURL, &abono.Notas, &abono.CreatedAt,
	)
	if err != nil {
		return models.Abono{}, err
	}

	// 4. Actualizar el estado del periodo
	// Obtenemos el total pagado hasta ahora en este periodo para permitir abonos parciales
	var totalPagado float64
	queryTotal := `SELECT COALESCE(SUM(monto_pagado), 0) FROM abonos WHERE periodo_pago_id = $1`
	err = tx.QueryRow(context.Background(), queryTotal, input.PeriodoPagoID).Scan(&totalPagado)
	if err != nil {
		return models.Abono{}, err
	}

	nuevoEstado := "pendiente"
	if totalPagado > 0 {
		nuevoEstado = "pagado"
	}

	queryUpdatePeriodo := `
		UPDATE periodos_pago 
		SET estado = $1, mora_aplicada = $2 
		WHERE id = $3
	`
	_, err = tx.Exec(context.Background(), queryUpdatePeriodo, nuevoEstado, moraAplicada, input.PeriodoPagoID)
	if err != nil {
		return models.Abono{}, err
	}

	// 5. Upgrade terrain to "vendido" if they have more than 1 abono
	var abonosCount int
	queryCount := `
		SELECT COUNT(*) 
		FROM abonos a
		JOIN periodos_pago pp ON a.periodo_pago_id = pp.id
		WHERE pp.plan_id = $1
	`
	err = tx.QueryRow(context.Background(), queryCount, periodo.PlanID).Scan(&abonosCount)
	if err == nil && abonosCount > 1 {
		queryUpgradeTerreno := `
			UPDATE terrenos t
			SET estado = 'vendido'
			FROM planes_pago p
			WHERE p.id = $1 AND p.terreno_id = t.id AND t.estado = 'apartado'
		`
		tx.Exec(context.Background(), queryUpgradeTerreno, periodo.PlanID)
	}

	// 6. Registrar comisión de ventas en egresos automáticamente (6.75%)
	comision := input.MontoPagado * 0.0675
	queryEgreso := `
		INSERT INTO egresos (parcela_id, fecha, concepto, monto, categoria, descripcion)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(
		context.Background(),
		queryEgreso,
		parcelaID,
		input.FechaPago,
		"Comisión a vendedores (6.75%)",
		comision,
		"Comisiones",
		"Comisión automática del 6.75% sobre el abono/venta registrado",
	)
	if err != nil {
		return models.Abono{}, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return models.Abono{}, err
	}

	return abono, nil
}

func GetAbonosByPeriodo(periodoID string) ([]models.Abono, error) {
	query := `
		SELECT 
			a.id, a.parcela_id, a.periodo_pago_id, a.numero_abono, a.monto_pagado, a.moneda, a.tipo_cambio, a.fecha_pago, a.metodo_pago, a.comprobante_url, a.notas, a.created_at,
			COALESCE(t.clave, '') as terreno_clave,
			COALESCE(t.nombre, '') as terreno_nombre,
			COALESCE(c.nombre_completo, '') as cliente_nombre,
			t.id as terreno_id
		FROM abonos a
		LEFT JOIN periodos_pago pp ON a.periodo_pago_id = pp.id
		LEFT JOIN planes_pago plan ON pp.plan_id = plan.id
		LEFT JOIN terrenos t ON plan.terreno_id = t.id
		LEFT JOIN clientes c ON plan.cliente_id = c.id
		WHERE a.periodo_pago_id = $1 
		ORDER BY a.created_at DESC
	`
	rows, err := database.DB.Query(context.Background(), query, periodoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var abonos []models.Abono
	for rows.Next() {
		var a models.Abono
		if err := rows.Scan(
			&a.ID, &a.ParcelaID, &a.PeriodoPagoID, &a.NumeroAbono, &a.MontoPagado, 
			&a.Moneda, &a.TipoCambio, &a.FechaPago, &a.MetodoPago, &a.ComprobanteURL, &a.Notas, &a.CreatedAt,
			&a.TerrenoClave, &a.TerrenoNombre, &a.ClienteNombre, &a.TerrenoID,
		); err != nil {
			return nil, err
		}
		abonos = append(abonos, a)
	}
	return abonos, nil
}

func GetAllAbonos(parcelaID string) ([]models.Abono, error) {
	query := `
		SELECT 
			a.id, a.parcela_id, a.periodo_pago_id, a.numero_abono, a.monto_pagado, a.moneda, a.tipo_cambio, a.fecha_pago, a.metodo_pago, a.comprobante_url, a.notas, a.created_at,
			COALESCE(t.clave, '') as terreno_clave,
			COALESCE(t.nombre, '') as terreno_nombre,
			COALESCE(c.nombre_completo, '') as cliente_nombre,
			t.id as terreno_id
		FROM abonos a
		LEFT JOIN periodos_pago pp ON a.periodo_pago_id = pp.id
		LEFT JOIN planes_pago plan ON pp.plan_id = plan.id
		LEFT JOIN terrenos t ON plan.terreno_id = t.id
		LEFT JOIN clientes c ON plan.cliente_id = c.id
		WHERE a.parcela_id = $1 
		ORDER BY a.fecha_pago DESC, a.created_at DESC
	`
	rows, err := database.DB.Query(context.Background(), query, parcelaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var abonos []models.Abono
	for rows.Next() {
		var a models.Abono
		if err := rows.Scan(
			&a.ID, &a.ParcelaID, &a.PeriodoPagoID, &a.NumeroAbono, &a.MontoPagado, 
			&a.Moneda, &a.TipoCambio, &a.FechaPago, &a.MetodoPago, &a.ComprobanteURL, &a.Notas, &a.CreatedAt,
			&a.TerrenoClave, &a.TerrenoNombre, &a.ClienteNombre, &a.TerrenoID,
		); err != nil {
			return nil, err
		}
		abonos = append(abonos, a)
	}
	return abonos, nil
}

func UpdateAbonoComprobante(abonoID string, comprobanteURL string) error {
	query := `
		UPDATE abonos
		SET comprobante_url = $1
		WHERE id = $2
	`
	_, err := database.DB.Exec(context.Background(), query, comprobanteURL, abonoID)
	return err
}

func UpdateAbono(abonoID string, input models.AbonoInput) (models.Abono, error) {
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return models.Abono{}, err
	}
	defer tx.Rollback(context.Background())

	queryAbono := `
		UPDATE abonos
		SET monto_pagado = $1, moneda = $2, tipo_cambio = $3, fecha_pago = $4, metodo_pago = $5, notas = $6
		WHERE id = $7
		RETURNING id, parcela_id, periodo_pago_id, numero_abono, monto_pagado, moneda, tipo_cambio, fecha_pago, metodo_pago, comprobante_url, notas, created_at
	`
	var abono models.Abono
	var tipoCambio *float64
	if input.TipoCambio > 0 {
		tipoCambio = &input.TipoCambio
	}
	err = tx.QueryRow(
		context.Background(), queryAbono,
		input.MontoPagado, input.Moneda, tipoCambio, input.FechaPago, input.MetodoPago, input.Notas, abonoID,
	).Scan(
		&abono.ID, &abono.ParcelaID, &abono.PeriodoPagoID, &abono.NumeroAbono, &abono.MontoPagado, 
		&abono.Moneda, &abono.TipoCambio, &abono.FechaPago, &abono.MetodoPago, &abono.ComprobanteURL, &abono.Notas, &abono.CreatedAt,
	)
	if err != nil {
		return models.Abono{}, err
	}

	// Actualizar estado del periodo
	var totalPagado float64
	queryTotal := `SELECT COALESCE(SUM(monto_pagado), 0) FROM abonos WHERE periodo_pago_id = $1`
	err = tx.QueryRow(context.Background(), queryTotal, abono.PeriodoPagoID).Scan(&totalPagado)
	if err != nil {
		return models.Abono{}, err
	}
	nuevoEstado := "pendiente"
	if totalPagado > 0 {
		nuevoEstado = "pagado"
	}
	queryUpdatePeriodo := `UPDATE periodos_pago SET estado = $1 WHERE id = $2`
	_, err = tx.Exec(context.Background(), queryUpdatePeriodo, nuevoEstado, abono.PeriodoPagoID)
	if err != nil {
		return models.Abono{}, err
	}

	err = tx.Commit(context.Background())
	return abono, err
}

func DeleteAbono(abonoID string) error {
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var periodoID string
	queryGet := `SELECT periodo_pago_id FROM abonos WHERE id = $1`
	err = tx.QueryRow(context.Background(), queryGet, abonoID).Scan(&periodoID)
	if err != nil {
		return err
	}

	queryDel := `DELETE FROM abonos WHERE id = $1`
	_, err = tx.Exec(context.Background(), queryDel, abonoID)
	if err != nil {
		return err
	}

	var totalPagado float64
	queryTotal := `SELECT COALESCE(SUM(monto_pagado), 0) FROM abonos WHERE periodo_pago_id = $1`
	err = tx.QueryRow(context.Background(), queryTotal, periodoID).Scan(&totalPagado)
	if err != nil {
		return err
	}

	nuevoEstado := "pendiente"
	if totalPagado > 0 {
		nuevoEstado = "pagado"
	}
	queryUpdatePeriodo := `UPDATE periodos_pago SET estado = $1 WHERE id = $2`
	_, err = tx.Exec(context.Background(), queryUpdatePeriodo, nuevoEstado, periodoID)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
