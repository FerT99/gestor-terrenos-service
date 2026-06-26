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

	query := `SELECT id, nombre, email, rol FROM usuarios`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, nombre, email, rol string
		if err := rows.Scan(&id, &nombre, &email, &rol); err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("ID: %s, Nombre: '%s', Email: '%s', Rol: '%s'\n", id, nombre, email, rol)
	}
}
