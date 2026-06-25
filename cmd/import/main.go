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

		// Limpiar símbolo de dólar y comas de miles para Valor Total (formato US/MX: 160,000.00)
		valorStr := strings.ReplaceAll(row[ColValorTotal], "$", "")
		valorStr = strings.ReplaceAll(valorStr, " ", "")
		valorStr = strings.ReplaceAll(valorStr, ",", "") // Quitar comas de miles
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
				SET superficie_m2 = $1, precio_lista = $2, propietario = $3, estado = 'vendido' 
				WHERE id = $4`,
				m2, valorTotal, clienteNombre, terrenoID)
			if err != nil {
				log.Printf("Error actualizando terreno %s: %v", claveTerreno, err)
				continue
			}
			terreno.ID = terrenoID
		}

		// C. Calcular plazos fijos
		// Nos indicaron que TODOS los terrenos se dividen entre 40 abonos (plazos = 40)
		plazos := 40

		maxCol := len(row)
		if maxCol > 46 {
			maxCol = 46 // Solo leer hasta Abono 40 (índice 45)
		}

		var plan models.PlanPago
		err = database.DB.QueryRow(context.Background(), "SELECT id FROM planes_pago WHERE terreno_id = $1 LIMIT 1", terreno.ID).Scan(&plan.ID)
		if err != nil {
			plan, err = repository.CreatePlanPago(models.PlanPagoInput{
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
		}

		// D. Registrar Abonos
		// Obtener la fila de fechas (que siempre está debajo de la fila del cliente)
		var fechasRow []string
		if i+1 < len(records) {
			nextRow := records[i+1]
			// Confirmamos que es una fila de fechas si el nombre del cliente está vacío
			if len(nextRow) > ColCliente && strings.TrimSpace(nextRow[ColCliente]) == "" {
				fechasRow = nextRow
			}
		}

		periodos, _ := repository.GetPeriodosByPlan(plan.ID)

		abonoIndex := 0
		for j := ColAbonosInicio; j < maxCol; j++ {
			abonoStr := strings.TrimSpace(row[j])
			if abonoStr == "" {
				continue
			}

			// Limpiar formato moneda (US/MX)
			abonoStr = strings.ReplaceAll(abonoStr, "$", "")
			abonoStr = strings.ReplaceAll(abonoStr, " ", "")
			abonoStr = strings.ReplaceAll(abonoStr, ",", "")
			montoAbono, _ := strconv.ParseFloat(abonoStr, 64)

			if montoAbono > 0 && abonoIndex < len(periodos) {
				periodoActual := periodos[abonoIndex]

				// Determinar la fecha real del abono (de la fila de abajo)
				fechaPagoStr := fechaInicio.AddDate(0, abonoIndex, 0).Format("2006-01-02") // Fallback a simulada
				if len(fechasRow) > j {
					fStr := strings.TrimSpace(fechasRow[j])
					if fStr != "" {
						if t, err := time.Parse("02/01/06", fStr); err == nil {
							fechaPagoStr = t.Format("2006-01-02")
						}
					}
				}

				// Solo procesar si la fecha es de Junio o Julio (meses 6 y 7)
				parsedFecha, _ := time.Parse("2006-01-02", fechaPagoStr)
				if parsedFecha.Month() != time.June && parsedFecha.Month() != time.July {
					abonoIndex++
					continue
				}

				// Crear el abono
				_, err = repository.CreateAbono(parcelaID, models.AbonoInput{
					PeriodoPagoID: periodoActual.ID,
					MontoPagado:   montoAbono,
					FechaPago:     fechaPagoStr,
					MetodoPago:    "transferencia",
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
