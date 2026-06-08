package repository

import (
	"context"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/jackc/pgx/v5"
	"fmt"
	"strconv"
)

// nullStr convierte un string vacío a nil para campos nullable de la BD
func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func GetNextClave(parcelaID string) (string, error) {
	query := `
		SELECT clave
		FROM terrenos
		WHERE parcela_id = $1 AND clave ~ '^T[0-9]+$'
		ORDER BY CAST(SUBSTRING(clave FROM 2) AS INTEGER) DESC
		LIMIT 1
	`
	var maxClave string
	err := database.DB.QueryRow(context.Background(), query, parcelaID).Scan(&maxClave)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "T1", nil
		}
		return "", err
	}

	numStr := maxClave[1:]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "T1", nil // fallback en caso extraño
	}

	return fmt.Sprintf("T%d", num+1), nil
}

func GetAllTerrenos(parcelaID string, vendedorID *string) ([]models.Terreno, error) {
	query := `
		SELECT id, parcela_id, clave, nombre, fase, superficie_m2, precio_lista,
		       propietario, estado, coordenadas, notas, vendedor_id, created_at
		FROM terrenos
		WHERE parcela_id = $1
		ORDER BY created_at DESC
	`
	var rows pgx.Rows
	var err error
	if vendedorID != nil {
		query = `
			SELECT id, parcela_id, clave, nombre, fase, superficie_m2, precio_lista,
			       propietario, estado, coordenadas, notas, vendedor_id, created_at
			FROM terrenos
			WHERE parcela_id = $1 AND vendedor_id = $2
			ORDER BY created_at DESC
		`
		rows, err = database.DB.Query(context.Background(), query, parcelaID, *vendedorID)
	} else {
		rows, err = database.DB.Query(context.Background(), query, parcelaID)
	}
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terrenos []models.Terreno
	for rows.Next() {
		var t models.Terreno
		err := rows.Scan(
			&t.ID, &t.ParcelaID, &t.Clave, &t.Nombre, &t.Fase,
			&t.SuperficieM2, &t.PrecioLista,
			&t.Propietario, &t.Estado, &t.Coordenadas,
			&t.Notas, &t.VendedorID, &t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		terrenos = append(terrenos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if terrenos == nil {
		terrenos = []models.Terreno{}
	}
	return terrenos, nil
}

func GetTerrenoByID(id string) (models.Terreno, error) {
	query := `
		SELECT id, parcela_id, clave, nombre, fase, superficie_m2, precio_lista,
		       propietario, estado, coordenadas, notas, vendedor_id, created_at
		FROM terrenos
		WHERE id = $1
	`
	var t models.Terreno
	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&t.ID, &t.ParcelaID, &t.Clave, &t.Nombre, &t.Fase,
		&t.SuperficieM2, &t.PrecioLista,
		&t.Propietario, &t.Estado, &t.Coordenadas,
		&t.Notas, &t.VendedorID, &t.CreatedAt,
	)
	return t, err
}

func CreateTerreno(input models.TerrenoInput) (models.Terreno, error) {
	query := `
		INSERT INTO terrenos
		  (parcela_id, clave, nombre, fase, superficie_m2, precio_lista, propietario, estado, coordenadas, notas, vendedor_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, parcela_id, clave, nombre, fase, superficie_m2, precio_lista,
		          propietario, estado, coordenadas, notas, vendedor_id, created_at
	`
	var t models.Terreno
	err := database.DB.QueryRow(context.Background(), query,
		input.ParcelaID,
		input.Clave,
		nullStr(input.Nombre),
		nullStr(input.Fase),
		input.SuperficieM2,
		input.PrecioLista,
		nullStr(input.Propietario),
		input.Estado,
		nullStr(input.Coordenadas),
		nullStr(input.Notas),
		nullStr(func() string {
			if input.VendedorID != nil {
				return *input.VendedorID
			}
			return ""
		}()),
	).Scan(
		&t.ID, &t.ParcelaID, &t.Clave, &t.Nombre, &t.Fase,
		&t.SuperficieM2, &t.PrecioLista,
		&t.Propietario, &t.Estado, &t.Coordenadas,
		&t.Notas, &t.VendedorID, &t.CreatedAt,
	)
	return t, err
}

func UpdateTerreno(id string, input models.TerrenoInput) (models.Terreno, error) {
	query := `
		UPDATE terrenos
		SET parcela_id=$1, clave=$2, nombre=$3, fase=$4, superficie_m2=$5, precio_lista=$6,
		    propietario=$7, estado=$8, coordenadas=$9, notas=$10, vendedor_id=$11
		WHERE id=$12
		RETURNING id, parcela_id, clave, nombre, fase, superficie_m2, precio_lista,
		          propietario, estado, coordenadas, notas, vendedor_id, created_at
	`
	var t models.Terreno
	err := database.DB.QueryRow(context.Background(), query,
		input.ParcelaID,
		input.Clave,
		nullStr(input.Nombre),
		nullStr(input.Fase),
		input.SuperficieM2,
		input.PrecioLista,
		nullStr(input.Propietario),
		input.Estado,
		nullStr(input.Coordenadas),
		nullStr(input.Notas),
		nullStr(func() string {
			if input.VendedorID != nil {
				return *input.VendedorID
			}
			return ""
		}()),
		id,
	).Scan(
		&t.ID, &t.ParcelaID, &t.Clave, &t.Nombre, &t.Fase,
		&t.SuperficieM2, &t.PrecioLista,
		&t.Propietario, &t.Estado, &t.Coordenadas,
		&t.Notas, &t.VendedorID, &t.CreatedAt,
	)
	return t, err
}

func DeleteTerreno(id string) error {
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// 1. Eliminar abonos asociados a los planes de pago de este terreno
	_, err = tx.Exec(context.Background(), `
		DELETE FROM abonos 
		WHERE periodo_pago_id IN (
			SELECT id FROM periodos_pago WHERE plan_id IN (
				SELECT id FROM planes_pago WHERE terreno_id = $1
			)
		)
	`, id)
	if err != nil {
		return err
	}

	// 2. Eliminar periodos_pago
	_, err = tx.Exec(context.Background(), `
		DELETE FROM periodos_pago 
		WHERE plan_id IN (SELECT id FROM planes_pago WHERE terreno_id = $1)
	`, id)
	if err != nil {
		return err
	}

	// 3. Eliminar planes_pago
	_, err = tx.Exec(context.Background(), `
		DELETE FROM planes_pago WHERE terreno_id = $1
	`, id)
	if err != nil {
		return err
	}

	// 4. Eliminar el terreno
	_, err = tx.Exec(context.Background(), `DELETE FROM terrenos WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
