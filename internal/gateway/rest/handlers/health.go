package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"net/http"
)

type HealthHandler struct {
}

func (h *HealthHandler) Health(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"env":     config.AppConfig.ServerConfig.Port,
		"version": config.AppConfig.ServerConfig.Version,
	})
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
