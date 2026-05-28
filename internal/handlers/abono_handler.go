package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func CreateAbono(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}

	var input models.AbonoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}

	if input.MontoPagado <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "El monto pagado debe ser mayor a 0",
		})
	}

	abono, err := repository.CreateAbono(parcelaID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al registrar abono: " + err.Error(),
		})
	}

	accion := "CREAR_ABONO"
	if input.PerdonarMora {
		accion = "CREAR_ABONO_MORA_CONDONADA"
	}

	go repository.LogAction(models.AuditLogInput{
		UsuarioNombre: "Administrador", // Hardcodeado por ahora
		Accion:        accion,
		EntidadTipo:   "abonos",
		EntidadID:     abono.ID,
		Detalles: map[string]interface{}{
			"monto_pagado":  abono.MontoPagado,
			"moneda":        abono.Moneda,
			"metodo_pago":   input.MetodoPago,
			"perdonar_mora": input.PerdonarMora,
		},
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "data": abono})
}

func GetAbonos(c *fiber.Ctx) error {
	periodoID := c.Params("periodo_id")
	if periodoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID del periodo requerido",
		})
	}

	abonos, err := repository.GetAbonosByPeriodo(periodoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener abonos: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "data": abonos})
}

func GetAllAbonos(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}

	abonos, err := repository.GetAllAbonos(parcelaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener abonos: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "data": abonos})
}
