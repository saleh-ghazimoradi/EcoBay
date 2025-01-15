package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
	authService helper.Auth
}

func (u *UserHandler) Register(ctx *fiber.Ctx) error {
	user := dto.UserSignUp{}
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid user",
		})
	}

	token, err := u.userService.Signup(context.Background(), user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error signing up",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "register",
		"token":   token,
	})
}

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	loginInput := dto.UserLogin{}
	if err := ctx.BodyParser(&loginInput); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid user",
		})
	}

	token, err := u.userService.Login(context.Background(), loginInput.Email, loginInput.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "please provide correct user id or password",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "login",
		"token":   token,
	})
}

func (u *UserHandler) Verify(ctx *fiber.Ctx) error {
	user, err := u.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "please provide a valid user",
		})
	}

	var req dto.VerificationCodeInput
	if err = ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	if err = u.userService.VerifyCode(context.Background(), user.ID, req.Code); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verified successfully",
	})
}

func (u *UserHandler) GetVerificationCode(ctx *fiber.Ctx) error {
	user, err := u.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "please provide a valid user",
		})
	}

	code, err := u.userService.GetVerificationCode(context.Background(), user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to generate verification code",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
		"data":    code,
	})
}

func (u *UserHandler) CreateProfile(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}
func (u *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	user, err := u.authService.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error getting current user",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "get profile",
		"user":    user,
	})
}

func (u *UserHandler) AddToCart(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}

func (u *UserHandler) GetCart(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}

func (u *UserHandler) CreateOrder(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}

func (u *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}

func (u *UserHandler) GetOrder(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}

func (u *UserHandler) BecomeSeller(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
