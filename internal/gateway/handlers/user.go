package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"github.com/saleh-ghazimoradi/EcoBay/internal/validate"
)

type UserHandler struct {
	userService    service.UserService
	authentication helper.Authentication
}

func (u *UserHandler) Register(ctx *fiber.Ctx) error {
	payload := dto.UserSignup{}
	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	if err := validate.Validator.Struct(payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrValidate)
	}

	token, err := u.userService.SignUp(ctx.Context(), &payload)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrEmailExists):
			return errorResponse(ctx, fiber.StatusConflict, "a user with this email already exists")
		default:
			return serverErrorResponse(ctx, err)
		}
	}

	return successResponse(ctx, fiber.StatusCreated, "signed up successfully", token)
}

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	payload := dto.UserLogin{}
	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	if err := validate.Validator.Struct(payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrValidate)
	}

	token, err := u.userService.Login(ctx.Context(), &payload)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrUserNotFound):
			return notFoundResponse(ctx)
		case errors.Is(err, customErr.ErrInvalidPassword):
			return invalidCredentialsResponse(ctx)
		default:
			return serverErrorResponse(ctx, err)
		}
	}

	return successResponse(ctx, fiber.StatusCreated, "logged in successfully", token)
}

func (u *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	user := u.authentication.GetCurrentUser(ctx)
	return successResponse(ctx, fiber.StatusOK, "user profile", user)
}

func NewUserHandler(userService service.UserService, authentication helper.Authentication) *UserHandler {
	return &UserHandler{
		userService:    userService,
		authentication: authentication,
	}
}
