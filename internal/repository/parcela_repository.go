package repository

import (
	"context"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

func GetAllParcelas() ([]models.Parcela, error) {
	query := `
		SELECT id, nombre, created_at
		FROM parcelas
		ORDER BY created_at DESC
	`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcelas []models.Parcela
	for rows.Next() {
		var p models.Parcela
		err := rows.Scan(&p.ID, &p.Nombre, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcelas = append(parcelas, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if parcelas == nil {
		parcelas = []models.Parcela{}
	}
	return parcelas, nil
}

func GetParcelaByID(id string) (models.Parcela, error) {
	query := `
		SELECT id, nombre, created_at
		FROM parcelas
		WHERE id = $1
	`
	var p models.Parcela
	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&p.ID, &p.Nombre, &p.CreatedAt,
	)
	return p, err
}

func CreateParcela(input models.ParcelaInput) (models.Parcela, error) {
	query := `
		INSERT INTO parcelas (nombre)
		VALUES ($1)
		RETURNING id, nombre, created_at
	`
	var p models.Parcela
	err := database.DB.QueryRow(context.Background(), query,
		input.Nombre,
	).Scan(&p.ID, &p.Nombre, &p.CreatedAt)
	return p, err
}

func UpdateParcela(id string, input models.ParcelaInput) (models.Parcela, error) {
	query := `
		UPDATE parcelas
		SET nombre = $1
		WHERE id = $2
		RETURNING id, nombre, created_at
	`
	var p models.Parcela
	err := database.DB.QueryRow(context.Background(), query,
		input.Nombre, id,
	).Scan(&p.ID, &p.Nombre, &p.CreatedAt)
	return p, err
}

func DeleteParcela(id string) error {
	_, err := database.DB.Exec(context.Background(),
		`DELETE FROM parcelas WHERE id = $1`, id)
	return err
}
