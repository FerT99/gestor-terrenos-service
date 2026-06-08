package main

import (
	"context"
	"log"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Aviso: No se pudo cargar el archivo .env: %v", err)
	}

	// Conectar a la base de datos
	database.ConnectDB()
	defer database.DB.Close()

	log.Println("Iniciando limpieza total de base de datos para re-importación...")

	ctx := context.Background()
	
	// Limpiar datos en cascada usando TRUNCATE
	// Esto borra todos los abonos, periodos, planes, clientes y terrenos.
	query := `TRUNCATE TABLE abonos, periodos_pago, planes_pago, terrenos, clientes RESTART IDENTITY CASCADE;`
	
	_, err = database.DB.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Error limpiando la base de datos: %v", err)
	}

	log.Println("✅ Base de datos limpiada con éxito. Ya puedes correr el main.go para una importación limpia.")
}
