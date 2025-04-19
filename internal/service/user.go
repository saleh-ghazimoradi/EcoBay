package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"time"
)

type UserService interface {
	SignUp(ctx context.Context, input *dto.UserSignup) (string, error)
	Login(ctx context.Context, input *dto.UserLogin) (string, error)
	GetVerificationCode(ctx context.Context, input *domain.User) (string, error)
	VerifyCode(ctx context.Context, id uint, code string) error
}

type userService struct {
	userRepository repository.UserRepository
	authentication helper.Authentication
}

func (u *userService) findUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.userRepository.FindUserByEmail(ctx, email)
}

func (u *userService) SignUp(ctx context.Context, input *dto.UserSignup) (string, error) {
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
	user, err := u.findUserByEmail(ctx, input.Email)
	if err != nil {
		return "", customErr.ErrUserNotFound
	}

	if err = u.authentication.VerifyPassword(input.Password, user.Password); err != nil {
		return "", customErr.ErrInvalidPassword
	}

	return u.authentication.GenerateToken(user.ID, user.Email, user.UserType)
}

func (u *userService) findUserById(ctx context.Context, id uint) (*domain.User, error) {
	return u.userRepository.FindUserById(ctx, id)
}

func (u *userService) isVerifiedUser(ctx context.Context, id uint) bool {
	user, err := u.findUserById(ctx, id)
	return err == nil && user.Verified
}

func (u *userService) GetVerificationCode(ctx context.Context, input *domain.User) (string, error) {
	if u.isVerifiedUser(ctx, input.ID) {
		return "", errors.New("user is already verified")
	}

	code, err := u.authentication.GenerateCode()
	if err != nil {
		return "", err
	}

	user := &domain.User{
		Code:   code,
		Expiry: time.Now().Add(config.AppConfig.AuthConfig.CodeExpiry),
	}

	_, err = u.userRepository.UpdateUser(ctx, input.ID, user)
	if err != nil {
		return "", errors.New("unable to update verification code")
	}

	return code, nil
}

func (u *userService) VerifyCode(ctx context.Context, id uint, code string) error {
	if u.isVerifiedUser(ctx, id) {
		return errors.New("user is already verified")
	}

	user, err := u.findUserById(ctx, id)
	if err != nil {
		return err
	}

	if user.Code != code {
		return errors.New("invalid code")
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("code expired")
	}

	updateUser := &domain.User{
		Verified: true,
	}

	_, err = u.userRepository.UpdateUser(ctx, id, updateUser)
	if err != nil {
		return errors.New("unable to verify user")
	}

	return nil
}

func NewUserService(userRepository repository.UserRepository, authentication helper.Authentication) UserService {
	return &userService{
		userRepository: userRepository,
		authentication: authentication,
	}
}
