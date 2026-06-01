package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No env")
	}

	database.ConnectDB()
	defer database.DB.Close()

	// Buscar clientes basura
	rows, err := database.DB.Query(context.Background(), "SELECT id, nombre_completo FROM clientes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var idsToDelete []string
	for rows.Next() {
		var id, nombre string
		if err := rows.Scan(&id, &nombre); err != nil {
			continue
		}
		
		// Si el nombre contiene $ o /, o empieza con número, es basura de la importación
		if strings.Contains(nombre, "$") || strings.Contains(nombre, "/") || (len(nombre) > 0 && nombre[0] >= '0' && nombre[0] <= '9') {
			idsToDelete = append(idsToDelete, id)
		}
	}

	for _, id := range idsToDelete {
		// Usa la función de cascada que creamos antes para limpiar todo lo relacionado al cliente
		err := repository.DeleteCliente(id)
		if err != nil {
			log.Printf("Error borrando cliente %s: %v", id, err)
		}
	}

	// Borrar terrenos basura (los LOTE-IMP-%)
	_, err = database.DB.Exec(context.Background(), "DELETE FROM terrenos WHERE clave LIKE 'LOTE-IMP-%'")
	if err != nil {
		log.Printf("Error borrando terrenos: %v", err)
	}

	fmt.Printf("¡Limpieza completa! Se borraron %d clientes basura y todos los terrenos LOTE-IMP.\n", len(idsToDelete))
}
