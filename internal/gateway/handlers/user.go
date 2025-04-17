package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
)

type UserHandler struct {
	userService    service.UserService
	authentication helper.Authentication
}

func (u *UserHandler) Register(ctx *fiber.Ctx) error {
	user := dto.UserSignup{}
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid input",
		})
	}

	if user.Email == "" || user.Password == "" || user.Phone == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "email, password, and phone are required",
		})
	}

	token, err := u.userService.SignUp(ctx.UserContext(), &user)
	if err != nil {
		if err.Error() == "email already exists" {
			return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"message": "email already exists",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error signing up",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"token":   token,
		"message": "successfully signed up",
	})
}

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	payload := dto.UserLogin{}
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide valid input",
		})
	}

	token, err := u.userService.Login(ctx.UserContext(), &payload)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "password does not match" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"message": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error logging in",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"token":   token,
		"message": "successfully logged in",
	})
}

func (u *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	user := u.authentication.GetCurrentUser(ctx)
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "get profile",
		"user":    user,
	})
}

func NewUserHandler(userService service.UserService, authentication helper.Authentication) *UserHandler {
	return &UserHandler{
		userService:    userService,
		authentication: authentication,
	}
}
