package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/joho/godotenv"
)

// Configuraciones de columnas (índices basados en 0)
// AJUSTA ESTOS ÍNDICES SEGÚN TU EXCEL EXPORTADO A CSV
const (
	ColCliente      = 2 // Columna C (Cliente)
	ColClaveTerreno = 1 // Columna B (Número de Lote)
	ColM2           = 3 // Columna D (M2)
	ColValorTotal   = 4 // Columna E (Valor total)
	ColFechaInicio  = 5 // Columna F (Fecha de inicio)
	ColAbonosInicio = 6 // Columna G en adelante (Abonos)
)

func main() {
	// 1. Cargar variables de entorno y conectar a BD
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se encontró .env, usando variables de entorno del sistema")
	}

	database.ConnectDB()
	defer database.DB.Close()

	// Obtener la parcela principal (Ajusta el nombre si es diferente)
	var parcelaID string
	err = database.DB.QueryRow(context.Background(), "SELECT id FROM parcelas WHERE nombre = 'Parcela Principal' LIMIT 1").Scan(&parcelaID)
	if err != nil {
		log.Fatalf("Error obteniendo parcela: %v", err)
	}

	// 2. Leer archivo CSV
	filePath := "datos.csv" // Pon tu archivo CSV aquí
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("No se pudo abrir el archivo %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// reader.Comma = ';' // Descomenta si tu Excel usa punto y coma

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error leyendo CSV: %v", err)
	}

	// 3. Procesar filas (saltando la cabecera)
	for i, row := range records {
		if i == 0 {
			continue // Saltar encabezados
		}

		clienteNombre := strings.TrimSpace(row[ColCliente])
		if clienteNombre == "" {
			continue
		}

		// Parsear datos
		m2, _ := strconv.ParseFloat(strings.ReplaceAll(row[ColM2], ",", "."), 64)

		// Limpiar símbolo de dólar y puntos de miles para Valor Total
		valorStr := strings.ReplaceAll(row[ColValorTotal], "$", "")
		valorStr = strings.ReplaceAll(valorStr, ".", "")
		valorStr = strings.ReplaceAll(valorStr, ",", ".") // Si los decimales usan coma
		valorTotal, _ := strconv.ParseFloat(strings.TrimSpace(valorStr), 64)

		// Fecha de inicio (Asumiendo formato DD/MM/YY)
		fechaInicio, err := time.Parse("02/01/06", strings.TrimSpace(row[ColFechaInicio]))
		if err != nil {
			log.Printf("Error parseando fecha fila %d: %v", i, err)
			fechaInicio = time.Now()
		}

		// A. Crear o buscar Cliente
		var clienteID string
		err = database.DB.QueryRow(context.Background(), "SELECT id FROM clientes WHERE nombre_completo = $1 LIMIT 1", clienteNombre).Scan(&clienteID)
		if err != nil {
			// Crear cliente
			nuevoCliente, err := repository.CreateCliente(models.ClienteInput{
				NombreCompleto: clienteNombre,
				Estado:         "Activo",
				ParcelaID:      parcelaID,
			})
			if err != nil {
				log.Printf("Error creando cliente %s: %v", clienteNombre, err)
				continue
			}
			clienteID = nuevoCliente.ID
		}

		// B. Crear o Buscar Terreno
		claveTerreno := fmt.Sprintf("LOTE-IMP-%d", i) // Autogenerado por defecto
		if len(row) > ColClaveTerreno && row[ColClaveTerreno] != "" {
			claveTerreno = strings.TrimSpace(row[ColClaveTerreno])
		}

		var terreno models.Terreno
		var terrenoID string
		err = database.DB.QueryRow(context.Background(), "SELECT id FROM terrenos WHERE clave = $1 LIMIT 1", claveTerreno).Scan(&terrenoID)
		
		if err != nil {
			// No existe, crearlo
			terreno, err = repository.CreateTerreno(models.TerrenoInput{
				ParcelaID:    parcelaID,
				Clave:        claveTerreno,
				SuperficieM2: m2,
				PrecioLista:  valorTotal,
				Precio:       valorTotal,
				Estado:       "vendido",
				Propietario:  clienteNombre,
				Moneda:       "MXN",
			})
			if err != nil {
				log.Printf("Error creando terreno para %s: %v", clienteNombre, err)
				continue
			}
		} else {
			// Ya existe, actualizarlo (para corregir si en el intento anterior se guardó mal)
			_, err = database.DB.Exec(context.Background(), `
				UPDATE terrenos 
				SET superficie_m2 = $1, precio_lista = $2, precio = $3, propietario = $4, estado = 'vendido' 
				WHERE id = $5`, 
				m2, valorTotal, valorTotal, clienteNombre, terrenoID)
			if err != nil {
				log.Printf("Error actualizando terreno %s: %v", claveTerreno, err)
				continue
			}
			terreno.ID = terrenoID
		}

		// C. Crear Plan de Pago (sin plazos inicialmente, o calculados por los abonos)
		plazos := 0
		for j := ColAbonosInicio; j < len(row); j++ {
			if strings.TrimSpace(row[j]) != "" {
				plazos++
			}
		}

		plan, err := repository.CreatePlanPago(models.PlanPagoInput{
			ParcelaID:   parcelaID,
			TerrenoID:   terreno.ID,
			ClienteID:   clienteID,
			MontoTotal:  valorTotal,
			Enganche:    0,
			Plazos:      plazos,
			TasaInteres: 0,
			FechaInicio: fechaInicio,
		})
		if err != nil {
			log.Printf("Error creando plan de pago: %v", err)
			continue
		}

		// D. Registrar Abonos
		periodos, _ := repository.GetPeriodosByPlan(plan.ID)

		abonoIndex := 0
		for j := ColAbonosInicio; j < len(row); j++ {
			abonoStr := strings.TrimSpace(row[j])
			if abonoStr == "" {
				continue
			}

			// Limpiar formato moneda
			abonoStr = strings.ReplaceAll(abonoStr, "$", "")
			abonoStr = strings.ReplaceAll(abonoStr, ".", "")
			montoAbono, _ := strconv.ParseFloat(abonoStr, 64)

			if montoAbono > 0 && abonoIndex < len(periodos) {
				periodoActual := periodos[abonoIndex]

				// Crear el abono
				_, err = repository.CreateAbono(parcelaID, models.AbonoInput{
					PeriodoPagoID: periodoActual.ID,
					MontoPagado:   montoAbono,
					FechaPago:     fechaInicio.AddDate(0, abonoIndex, 0).Format("2006-01-02"), // Fecha simulada
					MetodoPago:    "Transferencia (Importación)",
					Moneda:        "MXN",
				})

				if err == nil {
					// Marcar periodo como pagado
					database.DB.Exec(context.Background(), "UPDATE periodos_pago SET estado = 'pagado' WHERE id = $1", periodoActual.ID)
				}
				abonoIndex++
			}
		}

		fmt.Printf("Fila %d procesada con éxito: %s\n", i, clienteNombre)
	}

	fmt.Println("Importación finalizada.")
}
