package repository

import (
	"context"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

func GetAllUsuarios() ([]models.Usuario, error) {
	query := `
		SELECT id, email, nombre_completo, rol, created_at
		FROM usuarios
		ORDER BY nombre_completo ASC
	`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.ID, &u.Email, &u.NombreCompleto, &u.Rol, &u.CreatedAt); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if usuarios == nil {
		usuarios = []models.Usuario{}
	}
	return usuarios, nil
}

func CreateOrUpdateUsuario(input models.UsuarioInput) (models.Usuario, error) {
	query := `
		INSERT INTO usuarios (id, email, nombre_completo, rol)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE 
		SET email = EXCLUDED.email, 
		    nombre_completo = EXCLUDED.nombre_completo,
		    rol = EXCLUDED.rol
		RETURNING id, email, nombre_completo, rol, created_at
	`
	var u models.Usuario
	err := database.DB.QueryRow(context.Background(), query,
		input.ID, input.Email, input.NombreCompleto, input.Rol,
	).Scan(&u.ID, &u.Email, &u.NombreCompleto, &u.Rol, &u.CreatedAt)

	return u, err
}

func GetUsuarioByID(id string) (models.Usuario, error) {
	query := `
		SELECT id, email, nombre_completo, rol, created_at
		FROM usuarios
		WHERE id = $1
	`
	var u models.Usuario
	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&u.ID, &u.Email, &u.NombreCompleto, &u.Rol, &u.CreatedAt,
	)
	return u, err
}
