package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	database.ConnectDB()
	defer database.DB.Close()

	query := `SELECT metodo_pago, monto_pagado, fecha_pago FROM abonos WHERE metodo_pago ILIKE '%transf%'`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var metodo string
		var monto float64
		var fecha time.Time
		if err := rows.Scan(&metodo, &monto, &fecha); err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("Metodo: '%s', Monto: %.2f, Fecha: %s\n", metodo, monto, fecha.Format("2006-01-02"))
		count++
	}
	fmt.Printf("Total transferencias: %d\n", count)
}
