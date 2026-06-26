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

	query := `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'egresos'
	`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var col, typ string
		if err := rows.Scan(&col, &typ); err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("Column: %s, Type: %s\n", col, typ)
	}
}
