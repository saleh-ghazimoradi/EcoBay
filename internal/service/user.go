package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service/service_models"
	"github.com/saleh-ghazimoradi/EcoBay/pkg/notification"
	"time"
)

type UserService interface {
	Signup(ctx context.Context, input dto.UserSignUp) (string, error)
	findUserByEmail(ctx context.Context, email string) (*service_models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetVerificationCode(ctx context.Context, user *service_models.User) error
	VerifyCode(ctx context.Context, id uint, code int) error
	CreateProfile(ctx context.Context, id uint, input any) error
	GetProfile(ctx context.Context, id uint) (*service_models.User, error)
	UpdateProfile(ctx context.Context, id uint, input any) error
	//BecomeSeller(ctx context.Context, id uint, input dto.SellerInput) (string, error)
	FindCart(ctx context.Context, id uint) ([]any, error)
	CreateCart(ctx context.Context, input any, user service_models.User) ([]any, error)
	CreateOrder(ctx context.Context, user service_models.User) (int, error)
	GetOrders(ctx context.Context, user service_models.User) ([]any, error)
	GetOrderById(ctx context.Context, id uint, uId uint) (any, error)
	isVerifiedUser(ctx context.Context, id uint) bool
}

type userService struct {
	userRepository      repository.UserRepository
	authService         helper.Auth
	notificationService notification.NotificationsClient
}

func (u *userService) Signup(ctx context.Context, input dto.UserSignUp) (string, error) {
	hashedPassword, err := u.authService.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := u.userRepository.CreateUser(ctx, &service_models.User{
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
	})

	return u.authService.GenerateToken(user.ID, user.Email, user.UserType)
}

func (u *userService) findUserByEmail(ctx context.Context, email string) (*service_models.User, error) {
	return u.userRepository.FindUserByEmail(ctx, email)
}

func (u *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user does not exist with the provided email id")
	}
	if err = u.authService.VerifyPassword(password, user.Password); err != nil {
		return "", err
	}

	return u.authService.GenerateToken(user.ID, user.Email, user.UserType)
}

func (u *userService) GetVerificationCode(ctx context.Context, user *service_models.User) error {
	if u.isVerifiedUser(ctx, user.ID) {
		return errors.New("user already verified")
	}

	code, err := u.authService.GenerateCode()
	if err != nil {
		return err
	}

	us := &service_models.User{
		Expiry: time.Now().Add(config.AppConfig.Necessities.CodeExpiry),
		Code:   code,
	}

	_, err = u.userRepository.UpdateUser(ctx, user.ID, us)
	if err != nil {
		return errors.New("unable to update verification code")
	}

	us, _ = u.userRepository.FindUserByID(ctx, user.ID)

	msg := fmt.Sprintf("Your verification code is: %v", code)
	if err = u.notificationService.SendSMS(ctx, us.Phone, msg); err != nil {
		return errors.New("error on sending SMS")
	}

	return nil
}

func (u *userService) VerifyCode(ctx context.Context, id uint, code int) error {
	if u.isVerifiedUser(ctx, id) {
		return errors.New("user already verified")
	}

	user, err := u.userRepository.FindUserByID(ctx, id)
	if err != nil {
		return errors.New("user does not exist with the provided id")
	}

	if user.Code != code {
		return errors.New("verification code does not match")
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("verification code has expired")
	}

	updateUser := &service_models.User{
		Verified: true,
	}

	_, err = u.userRepository.UpdateUser(ctx, id, updateUser)
	if err != nil {
		return errors.New("unable to verify user")
	}

	return nil
}

func (u *userService) CreateProfile(ctx context.Context, id uint, input any) error {
	return nil
}

func (u *userService) GetProfile(ctx context.Context, id uint) (*service_models.User, error) {
	return nil, nil
}

func (u *userService) UpdateProfile(ctx context.Context, id uint, input any) error {
	return nil
}

//func (u *userService) BecomeSeller(ctx context.Context, id uint, input dto.SellerInput) (string, error) {
//	return "", nil
//}

func (u *userService) FindCart(ctx context.Context, id uint) ([]any, error) {
	return nil, nil
}

func (u *userService) CreateCart(ctx context.Context, input any, user service_models.User) ([]any, error) {
	return nil, nil
}

func (u *userService) CreateOrder(ctx context.Context, user service_models.User) (int, error) {
	return 0, nil
}

func (u *userService) GetOrders(ctx context.Context, user service_models.User) ([]any, error) {
	return nil, nil
}

func (u *userService) GetOrderById(ctx context.Context, id uint, uId uint) (any, error) {
	return nil, nil
}

func (u *userService) isVerifiedUser(ctx context.Context, id uint) bool {
	currentUser, err := u.userRepository.FindUserByID(ctx, id)
	return err == nil && currentUser.Verified
}

func NewUserService(userRepository repository.UserRepository, authService helper.Auth, notificationService notification.NotificationsClient) UserService {
	return &userService{
		userRepository:      userRepository,
		authService:         authService,
		notificationService: notificationService,
	}
}
