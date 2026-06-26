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

const (
	ColCliente      = 2 // Columna C (Cliente)
	ColClaveTerreno = 1 // Columna B (Número de Lote)
	ColM2           = 3 // Columna D (M2)
	ColValorTotal   = 4 // Columna E (Valor total)
	ColFechaInicio  = 5 // Columna F (Fecha de inicio)
	ColAbonosInicio = 6 // Columna G en adelante (Abonos)
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se encontró .env, usando variables de entorno del sistema")
	}

	database.ConnectDB()
	defer database.DB.Close()

	var parcelaID string
	err = database.DB.QueryRow(context.Background(), "SELECT id FROM parcelas WHERE nombre = 'Parcela Principal' LIMIT 1").Scan(&parcelaID)
	if err != nil {
		log.Fatalf("Error obteniendo parcela: %v", err)
	}

	filePath := "datos.csv"
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("No se pudo abrir el archivo %s: %v", filePath, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	// Some lines have different number of fields, let's not fail on it
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error leyendo CSV: %v", err)
	}

	for i, row := range records {
		if i == 0 || i == 1 {
			continue // Saltar encabezados
		}
		if len(row) <= ColCliente {
			continue
		}

		clienteNombre := strings.TrimSpace(row[ColCliente])
		if clienteNombre == "" {
			continue
		}

		// Rellenar la fila con strings vacíos hasta 46 columnas
		for len(row) < 46 {
			row = append(row, "")
		}

		m2Str := strings.ReplaceAll(row[ColM2], ".", "") // if there are thousands separator
		m2Str = strings.ReplaceAll(m2Str, ",", ".")
		m2, _ := strconv.ParseFloat(m2Str, 64)

		valorStr := strings.ReplaceAll(row[ColValorTotal], "$", "")
		valorStr = strings.ReplaceAll(valorStr, " ", "")
		valorStr = strings.ReplaceAll(valorStr, ".", "") // Quitar comas de miles que en español son puntos (160.000,00)
		valorStr = strings.ReplaceAll(valorStr, ",", ".") // Convertir coma decimal a punto (160000.00)
		valorTotal, _ := strconv.ParseFloat(strings.TrimSpace(valorStr), 64)

		fechaInicio, err := time.Parse("02/01/06", strings.TrimSpace(row[ColFechaInicio]))
		if err != nil {
			fechaInicio = time.Now()
		}

		var clienteID string
		err = database.DB.QueryRow(context.Background(), "SELECT id FROM clientes WHERE nombre_completo = $1 LIMIT 1", clienteNombre).Scan(&clienteID)
		if err != nil {
			nuevoCliente, err := repository.CreateCliente(models.ClienteInput{
				NombreCompleto: clienteNombre,
				Estado:         "Activo",
				ParcelaID:      parcelaID,
			})
			if err != nil {
				continue
			}
			clienteID = nuevoCliente.ID
		}

		claveTerreno := fmt.Sprintf("LOTE-IMP-%d", i)
		if len(row) > ColClaveTerreno && row[ColClaveTerreno] != "" {
			claveTerreno = strings.TrimSpace(row[ColClaveTerreno])
		}

		var terreno models.Terreno
		var terrenoID string
		err = database.DB.QueryRow(context.Background(), "SELECT id FROM terrenos WHERE clave = $1 LIMIT 1", claveTerreno).Scan(&terrenoID)

		if err != nil {
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
				continue
			}
		} else {
			// Update ownership just in case
			_, err = database.DB.Exec(context.Background(), `
				UPDATE terrenos 
				SET superficie_m2 = $1, precio_lista = $2, propietario = $3 
				WHERE id = $4`,
				m2, valorTotal, clienteNombre, terrenoID)
			terreno.ID = terrenoID
		}

		plazos := 40
		maxCol := len(row)
		if maxCol > 46 {
			maxCol = 46
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
				continue
			}
		}

		var fechasRow []string
		if i+1 < len(records) {
			nextRow := records[i+1]
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

			abonoStr = strings.ReplaceAll(abonoStr, "$", "")
			abonoStr = strings.ReplaceAll(abonoStr, " ", "")
			abonoStr = strings.ReplaceAll(abonoStr, ".", "")
			abonoStr = strings.ReplaceAll(abonoStr, ",", ".")
			montoAbono, _ := strconv.ParseFloat(abonoStr, 64)

			if montoAbono > 0 && abonoIndex < len(periodos) {
				periodoActual := periodos[abonoIndex]

				fechaPagoStr := fechaInicio.AddDate(0, abonoIndex, 0).Format("2006-01-02")
				if len(fechasRow) > j {
					fStr := strings.TrimSpace(fechasRow[j])
					if fStr != "" {
						if t, err := time.Parse("02/01/06", fStr); err == nil {
							fechaPagoStr = t.Format("2006-01-02")
						}
					}
				}

				metodoPago := "Efectivo"
				database.DB.Exec(context.Background(), "UPDATE abonos SET metodo_pago = $1 WHERE periodo_pago_id = $2", metodoPago, periodoActual.ID)

				_, err = repository.CreateAbono(parcelaID, models.AbonoInput{
					PeriodoPagoID: periodoActual.ID,
					MontoPagado:   montoAbono,
					FechaPago:     fechaPagoStr,
					MetodoPago:    metodoPago,
					Moneda:        "MXN",
				})

				if err == nil {
					database.DB.Exec(context.Background(), "UPDATE periodos_pago SET estado = 'pagado' WHERE id = $1", periodoActual.ID)
				}
				abonoIndex++
			}
		}
		fmt.Printf("Fila %d procesada con éxito: %s (Lote %s)\n", i, clienteNombre, claveTerreno)
	}
	fmt.Println("Importación CSV finalizada.")
}
