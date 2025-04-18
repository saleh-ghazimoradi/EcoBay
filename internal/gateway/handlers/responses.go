package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
)

func logError(ctx *fiber.Ctx, err error) {
	var (
		method = ctx.Method()
		uri    = ctx.OriginalURL()
	)

	slg.Logger.Error(err.Error(), "method", method, "uri", uri)
}

func errorResponse(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(&fiber.Map{
		"message": message,
	})
}

func successResponse(ctx *fiber.Ctx, status int, message string, data any) error {
	return ctx.Status(status).JSON(&fiber.Map{
		"message": message,
		"data":    data,
	})
}

func serverErrorResponse(ctx *fiber.Ctx, err error) error {
	logError(ctx, err)
	message := "the server encountered a problem and could not process your request"
	return errorResponse(ctx, fiber.StatusInternalServerError, message)
}

func badRequestResponse(ctx *fiber.Ctx, err error) error {
	return errorResponse(ctx, fiber.StatusBadRequest, err.Error())
}

func notFoundResponse(ctx *fiber.Ctx) error {
	message := "the requested resource could not be found"
	return errorResponse(ctx, fiber.StatusNotFound, message)
}

func failedValidationResponse(ctx *fiber.Ctx, errors map[string]string) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
		"errors": errors,
	})
}

func invalidCredentialsResponse(ctx *fiber.Ctx) error {
	message := "invalid authentication credentials"
	return errorResponse(ctx, fiber.StatusUnauthorized, message)
}

func rateLimitExceededResponse(ctx *fiber.Ctx) error {
	message := "rate limit exceeded"
	return errorResponse(ctx, fiber.StatusTooManyRequests, message)
}

func editConflictResponse(ctx *fiber.Ctx) error {
	message := "unable to update the record due to an edit conflict, please try again"
	return errorResponse(ctx, fiber.StatusConflict, message)
}

func invalidAuthenticationTokenResponse(ctx *fiber.Ctx) error {
	ctx.Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	return errorResponse(ctx, fiber.StatusUnauthorized, message)
}

func methodNotAllowedResponse(ctx *fiber.Ctx) error {
	message := fmt.Sprintf("the %s method is not supported for this resource", ctx.Method())
	return errorResponse(ctx, fiber.StatusMethodNotAllowed, message)
}

func authenticationRequiredResponse(ctx *fiber.Ctx) error {
	message := "you must be authenticated to access this resource"
	return errorResponse(ctx, fiber.StatusUnauthorized, message)
}

func notPermittedResponse(ctx *fiber.Ctx) error {
	message := "your user account does not have the necessary permissions to access this resource"
	return errorResponse(ctx, fiber.StatusForbidden, message)
}

func inactiveAccountResponse(ctx *fiber.Ctx) error {
	message := "your user account must be activated to access this resource"
	return errorResponse(ctx, fiber.StatusForbidden, message)
}
