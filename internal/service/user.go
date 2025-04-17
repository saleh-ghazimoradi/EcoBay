package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
)

type UserService interface {
	SignUp(ctx context.Context, input *dto.UserSignup) (string, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.userRepository.FindUserByEmail(ctx, email)
}

func (u *userService) SignUp(ctx context.Context, input *dto.UserSignup) (string, error) {
	if input.Email == "" || input.Password == "" || input.Phone == "" {
		slg.Logger.Error("email, password, and phone cannot be empty")
		return "", errors.New("email, password, and phone cannot be empty")
	}

	user, err := u.userRepository.CreateUser(ctx, &domain.User{
		Email:    input.Email,
		Password: input.Password,
		Phone:    input.Phone,
	})
	if err != nil {
		return "", err
	}

	userInfo := fmt.Sprintf("%v, %v,%v", user.ID, user.Email, user.UserType)
	return userInfo, nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
