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
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	rows, err := pool.Query(context.Background(), "SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = 'planes_pago'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Columnas de planes_pago:")
	for rows.Next() {
		var colName, dataType, isNullable string
		rows.Scan(&colName, &dataType, &isNullable)
		fmt.Printf("- %s (%s, nullable: %s)\n", colName, dataType, isNullable)
	}
}
