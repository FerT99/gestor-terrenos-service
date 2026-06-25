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
	
	query := `
		CREATE TABLE IF NOT EXISTS audit_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			usuario_nombre VARCHAR(255) NOT NULL,
			accion VARCHAR(255) NOT NULL,
			entidad_tipo VARCHAR(100) NOT NULL,
			entidad_id UUID,
			detalles JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := database.DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tabla audit_logs creada exitosamente.")
}
