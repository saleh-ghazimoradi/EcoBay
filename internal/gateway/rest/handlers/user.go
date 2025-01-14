package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
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
		"message": token,
	})
}

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "login",
	})
}

func (u *UserHandler) Verify(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verify",
	})
}

func (u *UserHandler) GetVerificationCode(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}

func (u *UserHandler) CreateProfile(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
	})
}
func (u *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verification code",
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
