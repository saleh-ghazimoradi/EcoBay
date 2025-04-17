package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
)

type UserService interface {
	SignUp(ctx context.Context, input *dto.UserSignup) (string, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	Login(ctx context.Context, input *dto.UserLogin) (string, error)
}

type userService struct {
	userRepository repository.UserRepository
	auth           helper.Auth
}

func (u *userService) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.userRepository.FindUserByEmail(ctx, email)
}

func (u *userService) SignUp(ctx context.Context, input *dto.UserSignup) (string, error) {
	if input.Email == "" || input.Password == "" || input.Phone == "" {
		slg.Logger.Error("email, password, and phone cannot be empty")
		return "", errors.New("email, password, and phone cannot be empty")
	}

	hashedPassword, err := u.auth.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := u.userRepository.CreateUser(ctx, &domain.User{
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
	})

	return u.auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (u *userService) Login(ctx context.Context, input *dto.UserLogin) (string, error) {
	user, err := u.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err = u.auth.VerifyPassword(input.Password, user.Password); err != nil {
		return "", err
	}

	return u.auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func NewUserService(userRepository repository.UserRepository, auth helper.Auth) UserService {
	return &userService{
		userRepository: userRepository,
		auth:           auth,
	}
}
