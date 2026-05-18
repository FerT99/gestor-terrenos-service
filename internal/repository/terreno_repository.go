package repository

import (
	"context"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

// nullStr convierte un string vacío a nil para campos nullable de la BD
func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func GetAllTerrenos() ([]models.Terreno, error) {
	query := `
		SELECT id, clave, nombre, fase, superficie_m2, precio_lista,
		       propietario, estado, coordenadas, notas, created_at
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
			&t.ID, &t.Clave, &t.Nombre, &t.Fase,
			&t.SuperficieM2, &t.PrecioLista,
			&t.Propietario, &t.Estado, &t.Coordenadas,
			&t.Notas, &t.CreatedAt,
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

func CreateTerreno(input models.TerrenoInput) (models.Terreno, error) {
	query := `
		INSERT INTO terrenos
		  (clave, nombre, fase, superficie_m2, precio_lista, propietario, estado, coordenadas, notas)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, clave, nombre, fase, superficie_m2, precio_lista,
		          propietario, estado, coordenadas, notas, created_at
	`
	var t models.Terreno
	err := database.DB.QueryRow(context.Background(), query,
		input.Clave,
		nullStr(input.Nombre),
		nullStr(input.Fase),
		input.SuperficieM2,
		input.PrecioLista,
		nullStr(input.Propietario),
		input.Estado,
		nullStr(input.Coordenadas),
		nullStr(input.Notas),
	).Scan(
		&t.ID, &t.Clave, &t.Nombre, &t.Fase,
		&t.SuperficieM2, &t.PrecioLista,
		&t.Propietario, &t.Estado, &t.Coordenadas,
		&t.Notas, &t.CreatedAt,
	)
	return t, err
}

func UpdateTerreno(id string, input models.TerrenoInput) (models.Terreno, error) {
	query := `
		UPDATE terrenos
		SET clave=$1, nombre=$2, fase=$3, superficie_m2=$4, precio_lista=$5,
		    propietario=$6, estado=$7, coordenadas=$8, notas=$9
		WHERE id=$10
		RETURNING id, clave, nombre, fase, superficie_m2, precio_lista,
		          propietario, estado, coordenadas, notas, created_at
	`
	var t models.Terreno
	err := database.DB.QueryRow(context.Background(), query,
		input.Clave,
		nullStr(input.Nombre),
		nullStr(input.Fase),
		input.SuperficieM2,
		input.PrecioLista,
		nullStr(input.Propietario),
		input.Estado,
		nullStr(input.Coordenadas),
		nullStr(input.Notas),
		id,
	).Scan(
		&t.ID, &t.Clave, &t.Nombre, &t.Fase,
		&t.SuperficieM2, &t.PrecioLista,
		&t.Propietario, &t.Estado, &t.Coordenadas,
		&t.Notas, &t.CreatedAt,
	)
	return t, err
}

func DeleteTerreno(id string) error {
	_, err := database.DB.Exec(context.Background(),
		`DELETE FROM terrenos WHERE id = $1`, id)
	return err
}
