package handlers

import "github.com/gofiber/fiber/v2"

type HealthCheckHandler struct {
}

func (h HealthCheckHandler) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "I'm breathing",
	})
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}
