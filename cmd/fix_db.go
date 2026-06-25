package main

import (
	"context"
	"fmt"
	"log"
	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No env file")
	}
	database.ConnectDB()
	defer database.DB.Close()

	// Ver abonos sospechosos
	rows, err := database.DB.Query(context.Background(), "SELECT id, monto_pagado FROM abonos WHERE monto_pagado < 1000")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id string
		var monto float64
		rows.Scan(&id, &monto)
		count++
	}
	fmt.Printf("Encontrados %d abonos con monto menor a 1000\n", count)

	// Update abonos
	res, err := database.DB.Exec(context.Background(), "UPDATE abonos SET monto_pagado = monto_pagado * 1000 WHERE monto_pagado < 1000 AND monto_pagado > 0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Abonos actualizados: %d\n", res.RowsAffected())

	// Update planes de pago monto_total
	res2, err := database.DB.Exec(context.Background(), "UPDATE planes_pago SET monto_total = monto_total * 1000 WHERE monto_total < 1000 AND monto_total > 0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Planes de pago actualizados: %d\n", res2.RowsAffected())

	// Update terrenos precio_lista and precio
	res3, err := database.DB.Exec(context.Background(), "UPDATE terrenos SET precio_lista = precio_lista * 1000 WHERE precio_lista < 1000 AND precio_lista > 0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Terrenos actualizados: %d\n", res3.RowsAffected())

	// Update periodos_pago monto_esperado
	res4, err := database.DB.Exec(context.Background(), "UPDATE periodos_pago SET monto_esperado = monto_esperado * 1000 WHERE monto_esperado < 1000 AND monto_esperado > 0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Periodos de pago actualizados: %d\n", res4.RowsAffected())
}
