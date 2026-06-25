package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("No DATABASE_URL")
	}

	config, _ := pgxpool.ParseConfig(dbUrl)
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	query := `
		SELECT 
			a.id, a.parcela_id, a.periodo_pago_id, a.numero_abono, a.monto_pagado, a.moneda, a.tipo_cambio, a.fecha_pago, a.metodo_pago, a.comprobante_url, a.notas, a.created_at,
			COALESCE(t.clave, '') as terreno_clave,
			COALESCE(t.nombre, '') as terreno_nombre,
			COALESCE(c.nombre_completo, '') as cliente_nombre
		FROM abonos a
		LEFT JOIN periodos_pago pp ON a.periodo_pago_id = pp.id
		LEFT JOIN planes_pago plan ON pp.plan_id = plan.id
		LEFT JOIN terrenos t ON plan.terreno_id = t.id
		LEFT JOIN clientes c ON plan.cliente_id = c.id
	`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Fatal("Query error: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, parcelaID, periodoPagoID, moneda, terrenoClave, terrenoNombre, clienteNombre string
		var numeroAbono int
		var montoPagado float64
		var tipoCambio *float64
		var fechaPago interface{}
		var metodoPago, comprobanteURL, notas *string
		var createdAt interface{}

		err := rows.Scan(
			&id, &parcelaID, &periodoPagoID, &numeroAbono, &montoPagado, 
			&moneda, &tipoCambio, &fechaPago, &metodoPago, &comprobanteURL, &notas, &createdAt,
			&terrenoClave, &terrenoNombre, &clienteNombre,
		)
		if err != nil {
			fmt.Println("SCAN ERROR on row:", id, err)
		}
	}
	
	if rows.Err() != nil {
		fmt.Println("ROWS ERROR:", rows.Err())
	}
	fmt.Println("Done scanning.")
}
