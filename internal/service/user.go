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
	authentication helper.Authentication
}

func (u *userService) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.userRepository.FindUserByEmail(ctx, email)
}

func (u *userService) SignUp(ctx context.Context, input *dto.UserSignup) (string, error) {
	if input.Email == "" || input.Password == "" || input.Phone == "" {
		slg.Logger.Error("email, password, and phone cannot be empty")
		return "", errors.New("email, password, and phone cannot be empty")
	}

	if _, err := u.FindUserByEmail(ctx, input.Email); err == nil {
		slg.Logger.Error("email already exists", "email", input.Email)
		return "", errors.New("email already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return "", err
	}

	hashedPassword, err := u.authentication.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := u.userRepository.CreateUser(ctx, &domain.User{
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
	})
	if err != nil {
		return "", err
	}

	return u.authentication.GenerateToken(user.ID, user.Email, user.UserType)
}

func (u *userService) Login(ctx context.Context, input *dto.UserLogin) (string, error) {
	user, err := u.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err = u.authentication.VerifyPassword(input.Password, user.Password); err != nil {
		return "", err
	}

	return u.authentication.GenerateToken(user.ID, user.Email, user.UserType)
}

func NewUserService(userRepository repository.UserRepository, authentication helper.Authentication) UserService {
	return &userService{
		userRepository: userRepository,
		authentication: authentication,
	}
}
