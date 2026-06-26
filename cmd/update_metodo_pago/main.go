package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno y conectar a BD
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se encontró .env, usando variables de entorno del sistema")
	}

	database.ConnectDB()
	defer database.DB.Close()

	// Actualizar todos los abonos a 'Efectivo'
	query := `UPDATE abonos SET metodo_pago = 'Efectivo'`
	
	result, err := database.DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Error actualizando abonos: %v", err)
	}

	fmt.Printf("¡Actualización exitosa! Se modificaron %d abonos.\n", result.RowsAffected())
}
