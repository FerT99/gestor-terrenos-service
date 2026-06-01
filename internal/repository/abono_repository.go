package repository

import (
	"context"
	"errors"
	"time"

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

	// 2. Calcular Mora si la fecha de pago es mayor a la de vencimiento
	fechaPago, err := time.Parse("2006-01-02", input.FechaPago)
	if err != nil {
		return models.Abono{}, errors.New("fecha de pago inválida")
	}

	moraAplicada := 0.0
	// Considerar mora si pasó al menos 1 día de la fecha de vencimiento
	// truncated to day
	fPagoTrunc := time.Date(fechaPago.Year(), fechaPago.Month(), fechaPago.Day(), 0, 0, 0, 0, time.UTC)
	fVencTrunc := time.Date(periodo.FechaVencimiento.Year(), periodo.FechaVencimiento.Month(), periodo.FechaVencimiento.Day(), 0, 0, 0, 0, time.UTC)

	if fPagoTrunc.After(fVencTrunc) && !input.PerdonarMora {
		moraAplicada = periodo.MontoEsperado * 0.15
	}

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
		INSERT INTO abonos (parcela_id, periodo_pago_id, numero_abono, monto_pagado, moneda, fecha_pago, metodo_pago, comprobante_url, notas)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, parcela_id, periodo_pago_id, numero_abono, monto_pagado, moneda, fecha_pago, metodo_pago, comprobante_url, notas, created_at
	`
	var abono models.Abono
	var compURL *string
	if input.ComprobanteURL != "" {
		compURL = &input.ComprobanteURL
	}
	err = tx.QueryRow(
		context.Background(),
		queryAbono,
		parcelaID,
		input.PeriodoPagoID,
		numeroAbono,
		input.MontoPagado,
		input.Moneda,
		fechaPago,
		input.MetodoPago,
		compURL,
		input.Notas,
	).Scan(
		&abono.ID, &abono.ParcelaID, &abono.PeriodoPagoID, &abono.NumeroAbono, &abono.MontoPagado, 
		&abono.Moneda, &abono.FechaPago, &abono.MetodoPago, &abono.ComprobanteURL, &abono.Notas, &abono.CreatedAt,
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
	montoTotalRequerido := periodo.MontoEsperado + moraAplicada

	if totalPagado >= montoTotalRequerido {
		nuevoEstado = "pagado"
	} else if totalPagado > 0 {
		nuevoEstado = "parcial"
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

	err = tx.Commit(context.Background())
	if err != nil {
		return models.Abono{}, err
	}

	return abono, nil
}

func GetAbonosByPeriodo(periodoID string) ([]models.Abono, error) {
	query := `
		SELECT 
			a.id, a.parcela_id, a.periodo_pago_id, a.numero_abono, a.monto_pagado, a.moneda, a.fecha_pago, a.metodo_pago, a.comprobante_url, a.notas, a.created_at,
			COALESCE(t.clave, '') as terreno_clave,
			COALESCE(t.nombre, '') as terreno_nombre,
			COALESCE(c.nombre_completo, '') as cliente_nombre
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
			&a.Moneda, &a.FechaPago, &a.MetodoPago, &a.ComprobanteURL, &a.Notas, &a.CreatedAt,
			&a.TerrenoClave, &a.TerrenoNombre, &a.ClienteNombre,
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
			a.id, a.parcela_id, a.periodo_pago_id, a.numero_abono, a.monto_pagado, a.moneda, a.fecha_pago, a.metodo_pago, a.comprobante_url, a.notas, a.created_at,
			COALESCE(t.clave, '') as terreno_clave,
			COALESCE(t.nombre, '') as terreno_nombre,
			COALESCE(c.nombre_completo, '') as cliente_nombre
		FROM abonos a
		LEFT JOIN periodos_pago pp ON a.periodo_pago_id = pp.id
		LEFT JOIN planes_pago plan ON pp.plan_id = plan.id
		LEFT JOIN terrenos t ON plan.terreno_id = t.id
		LEFT JOIN clientes c ON plan.cliente_id = c.id
		WHERE a.parcela_id = $1 
		ORDER BY a.created_at DESC
		LIMIT 50
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
			&a.Moneda, &a.FechaPago, &a.MetodoPago, &a.ComprobanteURL, &a.Notas, &a.CreatedAt,
			&a.TerrenoClave, &a.TerrenoNombre, &a.ClienteNombre,
		); err != nil {
			return nil, err
		}
		abonos = append(abonos, a)
	}
	return abonos, nil
}
