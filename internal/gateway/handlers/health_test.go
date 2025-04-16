package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler_HealthCheck(t *testing.T) {
	app := fiber.New()
	handler := NewHealthCheckHandler()
	app.Get("/health", handler.HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err, "should not return an error when making the request")
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "should return 200 OK")

	var responseBody map[string]string

	if err = json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		assert.NoError(t, err, "should parse response body without error")
	}

	assert.Equal(t, "I'm breathing", responseBody["message"], "should return correct message")
}
