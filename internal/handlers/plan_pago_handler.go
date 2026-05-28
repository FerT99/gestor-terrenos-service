package handlers

import (
	"github.com/FerT99/gestor-terrenos-service/internal/models"
	"github.com/FerT99/gestor-terrenos-service/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func CreatePlanPago(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}

	var input models.PlanPagoInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Datos inválidos: " + err.Error(),
		})
	}
	input.ParcelaID = parcelaID

	plan, err := repository.CreatePlanPago(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al crear plan de pago: " + err.Error(),
		})
	}

	go repository.LogAction(models.AuditLogInput{
		UsuarioNombre: "Administrador", // Hardcodeado por ahora
		Accion:        "NUEVA_VENTA",
		EntidadTipo:   "planes_pago",
		EntidadID:     plan.ID,
		Detalles: map[string]interface{}{
			"monto_total": plan.MontoTotal,
			"enganche":    plan.Enganche,
			"moneda":      plan.Moneda,
		},
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "data": plan})
}

func GetPlanesPago(c *fiber.Ctx) error {
	parcelaID := c.Get("X-Parcela-Id")
	if parcelaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Header X-Parcela-Id requerido",
		})
	}

	planes, err := repository.GetAllPlanesPago(parcelaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener planes: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "data": planes})
}

func GetPeriodosPlan(c *fiber.Ctx) error {
	planID := c.Params("id")
	if planID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "ID del plan requerido",
		})
	}

	periodos, err := repository.GetPeriodosByPlan(planID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Error al obtener periodos: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{"error": false, "data": periodos})
}
