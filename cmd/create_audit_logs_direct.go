package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbUrl := "postgresql://postgres:Javierbarajas123%2A@db.lyeywamxqwwliigimhgj.supabase.co:5432/postgres"
	config, _ := pgxpool.ParseConfig(dbUrl)
	config.MaxConns = 1
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Ping falló: %v", err)
	}
	fmt.Println("Conectado directamente a la DB.")

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
	_, err = pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Error creando tabla: ", err)
	}
	fmt.Println("Tabla audit_logs creada.")
	pool.Close()
}
