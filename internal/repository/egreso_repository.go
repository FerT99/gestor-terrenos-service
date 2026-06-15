package repository

import (
	"context"
	"time"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

func parseFecha(f string) time.Time {
	t, err := time.Parse(time.RFC3339, f)
	if err == nil {
		return t
	}
	t, err = time.Parse("2006-01-02", f)
	if err == nil {
		return t
	}
	return time.Now()
}

func GetEgresosByParcela(parcelaID string) ([]models.Egreso, error) {
	query := `
		SELECT id, parcela_id, fecha, concepto, monto, categoria, descripcion
		FROM egresos
		WHERE parcela_id = $1
		ORDER BY fecha DESC
	`
	rows, err := database.DB.Query(context.Background(), query, parcelaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var egresos []models.Egreso
	for rows.Next() {
		var e models.Egreso
		if err := rows.Scan(&e.ID, &e.ParcelaID, &e.Fecha, &e.Concepto, &e.Monto, &e.Categoria, &e.Descripcion); err != nil {
			return nil, err
		}
		egresos = append(egresos, e)
	}
	return egresos, nil
}

func CreateEgreso(parcelaID string, input models.EgresoInput) (models.Egreso, error) {
	fecha := parseFecha(input.Fecha)

	query := `
		INSERT INTO egresos (parcela_id, fecha, concepto, monto, categoria, descripcion)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, parcela_id, fecha, concepto, monto, categoria, descripcion
	`
	var e models.Egreso
	err := database.DB.QueryRow(
		context.Background(),
		query,
		parcelaID,
		fecha,
		input.Concepto,
		input.Monto,
		input.Categoria,
		input.Descripcion,
	).Scan(&e.ID, &e.ParcelaID, &e.Fecha, &e.Concepto, &e.Monto, &e.Categoria, &e.Descripcion)

	return e, err
}

func UpdateEgreso(id string, parcelaID string, input models.EgresoInput) (models.Egreso, error) {
	fecha := parseFecha(input.Fecha)

	query := `
		UPDATE egresos
		SET fecha = $1, concepto = $2, monto = $3, categoria = $4, descripcion = $5
		WHERE id = $6 AND parcela_id = $7
		RETURNING id, parcela_id, fecha, concepto, monto, categoria, descripcion
	`
	var e models.Egreso
	err := database.DB.QueryRow(
		context.Background(),
		query,
		fecha,
		input.Concepto,
		input.Monto,
		input.Categoria,
		input.Descripcion,
		id,
		parcelaID,
	).Scan(&e.ID, &e.ParcelaID, &e.Fecha, &e.Concepto, &e.Monto, &e.Categoria, &e.Descripcion)

	return e, err
}

func DeleteEgreso(id string, parcelaID string) error {
	query := `DELETE FROM egresos WHERE id = $1 AND parcela_id = $2`
	_, err := database.DB.Exec(context.Background(), query, id, parcelaID)
	return err
}
