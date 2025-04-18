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
	user := dto.UserSignup{}
	if err := ctx.BodyParser(&user); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	if err := validate.Validator.Struct(user); err != nil {
		return badRequestResponse(ctx, customErr.ErrValidate)
	}

	token, err := u.userService.SignUp(ctx.UserContext(), &user)
	if err != nil {
		if errors.Is(err, customErr.ErrEmailExists) {
			return badRequestResponse(ctx, customErr.ErrEmailExists)
		}
		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusCreated, "successfully signed up", token)
}

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	payload := dto.UserLogin{}
	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	if err := validate.Validator.Struct(payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	token, err := u.userService.Login(ctx.UserContext(), &payload)
	if err != nil {
		if errors.Is(err, customErr.ErrNotFound) || errors.Is(err, customErr.ErrInvalidPassword) {
			return badRequestResponse(ctx, customErr.ErrInvalidInput)
		}

		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusCreated, "successfully logged in", token)
}

func (u *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	user := u.authentication.GetCurrentUser(ctx)

	return successResponse(ctx, fiber.StatusOK, "successfully retrieved user", user)
}

func NewUserHandler(userService service.UserService, authentication helper.Authentication) *UserHandler {
	return &UserHandler{
		userService:    userService,
		authentication: authentication,
	}
}
