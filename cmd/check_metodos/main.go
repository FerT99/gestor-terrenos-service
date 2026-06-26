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

	rows, err := database.DB.Query(context.Background(), "SELECT DISTINCT metodo_pago FROM abonos")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var metodo string
		if err := rows.Scan(&metodo); err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("Metodo: '%s'\n", metodo)
	}
}
