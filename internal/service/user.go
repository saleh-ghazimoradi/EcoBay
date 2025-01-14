package service

import (
	"context"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service/service_models"
	"log"
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
	userRepository repository.UserRepository
}

func (u *userService) Signup(ctx context.Context, input dto.UserSignUp) (string, error) {
	log.Println(input)

	return "this is my token", nil
}

func (u *userService) findUserByEmail(ctx context.Context, email string) (*service_models.User, error) {
	return nil, nil
}

func (u *userService) Login(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

func (u *userService) GetVerificationCode(ctx context.Context, user *service_models.User) error {
	return nil
}

func (u *userService) VerifyCode(ctx context.Context, id uint, code int) error {
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
	return false
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
