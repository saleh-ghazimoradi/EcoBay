package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
)

type UserHandler struct {
	userService service.UserService
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

	token, err := u.userService.SignUp(&user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error signing up",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"token":   token,
		"message": "successfully signed up",
	})
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
