package repository

import (
	"context"
	"database/sql"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/FerT99/gestor-terrenos-service/internal/models"
)

func sqlToNullStr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func GetAllClientes(parcelaID string) ([]models.Cliente, error) {
	query := `SELECT id, parcela_id, nombre_completo, email, telefono, direccion, estado, created_at FROM clientes WHERE parcela_id = $1 ORDER BY created_at DESC`
	rows, err := database.DB.Query(context.Background(), query, parcelaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var c models.Cliente
		var email, telefono, direccion sql.NullString
		if err := rows.Scan(&c.ID, &c.ParcelaID, &c.NombreCompleto, &email, &telefono, &direccion, &c.Estado, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.Email = sqlToNullStr(email)
		c.Telefono = sqlToNullStr(telefono)
		c.Direccion = sqlToNullStr(direccion)
		clientes = append(clientes, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if clientes == nil {
		clientes = []models.Cliente{}
	}
	return clientes, nil
}

func GetClienteByID(id string) (models.Cliente, error) {
	query := `SELECT id, parcela_id, nombre_completo, email, telefono, direccion, estado, created_at FROM clientes WHERE id = $1`
	var c models.Cliente
	var email, telefono, direccion sql.NullString
	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&c.ID, &c.ParcelaID, &c.NombreCompleto, &email, &telefono, &direccion, &c.Estado, &c.CreatedAt,
	)
	c.Email = sqlToNullStr(email)
	c.Telefono = sqlToNullStr(telefono)
	c.Direccion = sqlToNullStr(direccion)
	return c, err
}

func CreateCliente(input models.ClienteInput) (models.Cliente, error) {
	query := `
		INSERT INTO clientes (parcela_id, nombre_completo, email, telefono, direccion, estado)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, parcela_id, nombre_completo, email, telefono, direccion, estado, created_at
	`

	estado := input.Estado
	if estado == "" {
		estado = "Activo"
	}

	var c models.Cliente
	var email, telefono, direccion sql.NullString
	err := database.DB.QueryRow(
		context.Background(),
		query,
		input.ParcelaID,
		input.NombreCompleto,
		nullStr(input.Email),
		nullStr(input.Telefono),
		nullStr(input.Direccion),
		estado,
	).Scan(&c.ID, &c.ParcelaID, &c.NombreCompleto, &email, &telefono, &direccion, &c.Estado, &c.CreatedAt)

	c.Email = sqlToNullStr(email)
	c.Telefono = sqlToNullStr(telefono)
	c.Direccion = sqlToNullStr(direccion)
	return c, err
}

func UpdateCliente(id string, input models.ClienteInput) (models.Cliente, error) {
	query := `
		UPDATE clientes
		SET parcela_id=$1, nombre_completo=$2, email=$3, telefono=$4, direccion=$5, estado=$6
		WHERE id=$7
		RETURNING id, parcela_id, nombre_completo, email, telefono, direccion, estado, created_at
	`

	estado := input.Estado
	if estado == "" {
		estado = "Activo"
	}

	var c models.Cliente
	var email, telefono, direccion sql.NullString
	err := database.DB.QueryRow(
		context.Background(),
		query,
		input.ParcelaID,
		input.NombreCompleto,
		nullStr(input.Email),
		nullStr(input.Telefono),
		nullStr(input.Direccion),
		estado,
		id,
	).Scan(&c.ID, &c.ParcelaID, &c.NombreCompleto, &email, &telefono, &direccion, &c.Estado, &c.CreatedAt)

	c.Email = sqlToNullStr(email)
	c.Telefono = sqlToNullStr(telefono)
	c.Direccion = sqlToNullStr(direccion)
	return c, err
}

func DeleteCliente(id string) error {
	_, err := database.DB.Exec(context.Background(), `DELETE FROM clientes WHERE id = $1`, id)
	return err
}
