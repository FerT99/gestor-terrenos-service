package repository

import (
	"context"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

func GetAllTerrenos() ([]models.Terreno, error) {
	query := `
		SELECT id, clave, nombre, superficie_m2, precio, propietario_familiar, estado, notas, created_at, updated_at
		FROM terrenos
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terrenos []models.Terreno
	for rows.Next() {
		var t models.Terreno
		err := rows.Scan(
			&t.ID, &t.Clave, &t.Nombre, &t.SuperficieM2, &t.Precio, 
			&t.PropietarioFamiliar, &t.Estado, &t.Notas, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		terrenos = append(terrenos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Si está vacío, devolvemos un slice vacío en lugar de nil para JSON []
	if terrenos == nil {
		terrenos = []models.Terreno{}
	}

	return terrenos, nil
}
