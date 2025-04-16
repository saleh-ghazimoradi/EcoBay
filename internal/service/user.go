package service

import (
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
)

type UserService interface {
	SignUp(input *dto.UserSignup) (string, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) SignUp(input *dto.UserSignup) (string, error) {
	return "this is the token", nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
