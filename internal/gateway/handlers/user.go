package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"github.com/saleh-ghazimoradi/EcoBay/internal/validate"
	"strconv"
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

func (u *UserHandler) CreateProfile(ctx *fiber.Ctx) error {
	payload := dto.Profile{}
	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	user := u.authentication.GetCurrentUser(ctx)

	if err := u.userService.CreateProfile(ctx.Context(), user.ID, &payload); err != nil {
		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusOK, "profile successfully created", nil)
}

func (u *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	user := u.authentication.GetCurrentUser(ctx)

	profile, err := u.userService.GetProfile(ctx.Context(), user.ID)
	if err != nil {
		return notFoundResponse(ctx)
	}

	return successResponse(ctx, fiber.StatusOK, "profile successfully retrieved", profile)
}

func (u *UserHandler) UpdateProfile(ctx *fiber.Ctx) error {
	payload := dto.Profile{}
	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	user := u.authentication.GetCurrentUser(ctx)

	if err := u.userService.UpdateProfile(ctx.Context(), user.ID, &payload); err != nil {
		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusOK, "profile updated successfully", nil)
}

func (u *UserHandler) GetVerificationCode(ctx *fiber.Ctx) error {
	user := u.authentication.GetCurrentUser(ctx)

	err := u.userService.GetVerificationCode(ctx.Context(), user)
	if err != nil {
		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusOK, "user verification code", nil)
}

func (u *UserHandler) Verify(ctx *fiber.Ctx) error {
	payload := dto.VerificationCode{}
	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	if err := validate.Validator.Struct(payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrValidateCode)
	}

	user := u.authentication.GetCurrentUser(ctx)

	if err := u.userService.VerifyCode(ctx.Context(), user.ID, payload.Code); err != nil {
		return badRequestResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusOK, "user verified successfully", nil)
}

func (u *UserHandler) BecomeSeller(ctx *fiber.Ctx) error {
	payload := dto.BecomeSellerInput{}
	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	if err := validate.Validator.Struct(payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	user := u.authentication.GetCurrentUser(ctx)

	token, err := u.userService.BecomeSeller(ctx.Context(), user.ID, &payload)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "failed to become seller",
		})
	}

	return successResponse(ctx, fiber.StatusOK, "user successfully became seller", token)
}

func (u *UserHandler) AddToCart(ctx *fiber.Ctx) error {
	var payload dto.Cart
	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, customErr.ErrInvalidInput)
	}

	user := u.authentication.GetCurrentUser(ctx)

	cartItems, err := u.userService.CreateCart(ctx.Context(), &payload, user)
	if err != nil {
		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusCreated, "successfully added to the cart", cartItems)
}

func (u *UserHandler) GetCart(ctx *fiber.Ctx) error {
	user := u.authentication.GetCurrentUser(ctx)

	cart, _, err := u.userService.FindCart(ctx.Context(), user.ID)
	if err != nil {
		return notFoundResponse(ctx)
	}

	return successResponse(ctx, fiber.StatusOK, "cart successfully retrieved", cart)
}

func (u *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	user := u.authentication.GetCurrentUser(ctx)

	orders, err := u.userService.GetOrders(ctx.Context(), user)
	if err != nil {
		return notFoundResponse(ctx)
	}

	return successResponse(ctx, fiber.StatusOK, "orders successfully retrieved", orders)
}

func (u *UserHandler) GetOrder(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	uIntId := uint(id)

	user := u.authentication.GetCurrentUser(ctx)

	order, err := u.userService.GetOrderById(ctx.Context(), uIntId, user.ID)
	if err != nil {
		return notFoundResponse(ctx)
	}

	return successResponse(ctx, fiber.StatusOK, "order successfully retrieved", order)
}

func NewUserHandler(userService service.UserService, authentication helper.Authentication) *UserHandler {
	return &UserHandler{
		userService:    userService,
		authentication: authentication,
	}
}
