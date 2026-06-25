package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL no está configurada en .env")
	}

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("Error configurando la base de datos: %v\n", err)
	}

	// Supabase Pool settings recommendations
	config.MaxConns = 3

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v\n", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Base de datos no responde (ping falló): %v\n", err)
	}

	fmt.Println("¡Conectado exitosamente a PostgreSQL (Supabase)!")
	DB = pool
}
