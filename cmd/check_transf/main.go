package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	database.ConnectDB()
	defer database.DB.Close()

	query := `SELECT id, metodo_pago, monto_pagado, fecha_pago FROM abonos WHERE metodo_pago = 'Transferencia' OR metodo_pago = 'transferencia'`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, metodo string
		var monto float64
		var fecha string
		if err := rows.Scan(&id, &metodo, &monto, &fecha); err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("ID: %s, Metodo: '%s', Monto: %.2f, Fecha: %s\n", id, metodo, monto, fecha)
		count++
	}
	fmt.Printf("Total found: %d\n", count)
}
