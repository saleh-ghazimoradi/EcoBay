package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/pkg/notification"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"gorm.io/gorm"
	"time"
)

type UserService interface {
	SignUp(ctx context.Context, input *dto.UserSignup) (string, error)
	Login(ctx context.Context, input *dto.UserLogin) (string, error)
	GetVerificationCode(ctx context.Context, input *domain.User) error
	VerifyCode(ctx context.Context, id uint, code string) error
	BecomeSeller(ctx context.Context, id uint, input *dto.BecomeSellerInput) (string, error)
	CreateCart(ctx context.Context, input *dto.Cart, user *domain.User) ([]*domain.Cart, error)
	FindCart(ctx context.Context, id uint) ([]*domain.Cart, error)
	CreateOrder(ctx context.Context, uId uint, orderRef string, pId string, amount float64) error
	CreateProfile(ctx context.Context, id uint, input *dto.Profile) error
	GetOrderById(ctx context.Context, id, uId uint) (*domain.Order, error)
	GetOrders(ctx context.Context, user *domain.User) ([]*domain.Order, error)
	GetProfile(ctx context.Context, id uint) (*domain.User, error)
	UpdateProfile(ctx context.Context, id uint, input *dto.Profile) error
}

type userService struct {
	userRepository    repository.UserRepository
	catalogRepository repository.CatalogRepository
	authentication    helper.Authentication
	email             notification.Email
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

func (u *userService) GetVerificationCode(ctx context.Context, input *domain.User) error {
	if u.isVerifiedUser(ctx, input.ID) {
		return errors.New("user is already verified")
	}

	code, err := u.authentication.GenerateCode()
	if err != nil {
		return err
	}

	user := &domain.User{
		Code:   code,
		Expiry: time.Now().Add(config.AppConfig.AuthConfig.CodeExpiry),
	}

	_, err = u.userRepository.UpdateUser(ctx, input.ID, user)
	if err != nil {
		return errors.New("unable to update verification code")
	}

	user, _ = u.findUserById(ctx, input.ID)

	msg := fmt.Sprintf("Your verification code is %s", code)
	fmt.Println(msg)

	background(func() {
		data := map[string]any{
			"code": code,
			"id":   user.ID,
		}

		err = u.email.SendEmail(user.Email, "user_verification_code.tmpl", data)
		if err != nil {
			slg.Logger.Error(err.Error())
		}
	})

	return nil
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

func (u *userService) BecomeSeller(ctx context.Context, id uint, input *dto.BecomeSellerInput) (string, error) {
	user, _ := u.findUserById(ctx, id)

	if user.UserType == domain.Seller {
		return "", errors.New("user is already seller")
	}

	seller, err := u.userRepository.UpdateUser(ctx, id, &domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.Seller,
	})

	if err != nil {
		return "", err
	}

	token, err := u.authentication.GenerateToken(user.ID, user.Email, seller.UserType)
	if err != nil {
		return "", err
	}

	account := &domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	}

	if err = u.userRepository.CreateBankAccount(ctx, account); err != nil {
		slg.Logger.Error("failed to create bank account", "user_id", id, "error", err)
		return "", errors.New("failed to create bank account")
	}

	return token, nil
}

func (u *userService) CreateCart(ctx context.Context, input *dto.Cart, user *domain.User) ([]*domain.Cart, error) {

	if input.ProductId == 0 {
		return nil, errors.New("invalid product id")
	}
	if input.Qty == 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	cart, err := u.userRepository.FindCartItem(ctx, user.ID, input.ProductId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slg.Logger.Error("error finding cart item", "user_id", user.ID, "product_id", input.ProductId, "err", err)
		return nil, errors.New("error checking cart item")
	}

	if cart.ID > 0 {
		// Cart item exists, update or delete
		if input.Qty < 1 {
			if err := u.userRepository.DeleteCartById(ctx, cart.ID); err != nil {
				slg.Logger.Error("error deleting cart item", "cart_id", cart.ID, "err", err)
				return nil, errors.New("error deleting cart item")
			}
		} else {
			cart.Qty = input.Qty
			if err := u.userRepository.UpdateCart(ctx, cart); err != nil {
				slg.Logger.Error("error updating cart item", "cart_id", cart.ID, "err", err)
				return nil, errors.New("error updating cart item")
			}
		}
	} else {
		// Cart item doesn't exist, create new
		product, err := u.catalogRepository.FindProductById(ctx, input.ProductId)
		if err != nil {
			if errors.Is(err, customErr.ErrNotFound) {
				slg.Logger.Error("product not found", "product_id", input.ProductId)
				return nil, errors.New("product not found")
			}
			slg.Logger.Error("error retrieving product", "product_id", input.ProductId, "err", err)
			return nil, errors.New("error retrieving product")
		}

		// Log product details for debugging
		slg.Logger.Info("product found", "product_id", product.ID, "name", product.Name, "stock", product.Stock)

		// Create new cart item
		newCart := &domain.Cart{
			ProductId: input.ProductId,
			UserId:    user.ID,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Qty:       input.Qty,
			Price:     product.Price,
			SellerId:  product.UserId,
		}

		if err := u.userRepository.CreateCart(ctx, newCart); err != nil {
			slg.Logger.Error("error creating cart item", "user_id", user.ID, "product_id", input.ProductId, "err", err)
			return nil, fmt.Errorf("error creating cart: %w", err)
		}
	}

	// Fetch updated cart items
	cartItems, err := u.userRepository.FindCartItems(ctx, user.ID)
	if err != nil {
		slg.Logger.Error("error fetching cart items", "user_id", user.ID, "err", err)
		return nil, errors.New("error retrieving cart items")
	}

	return cartItems, nil
}

func (u *userService) FindCart(ctx context.Context, id uint) ([]*domain.Cart, error) {
	cartItems, err := u.userRepository.FindCartItems(ctx, id)
	if err != nil {
		slg.Logger.Error("error on finding cart items", "err", err)
		return nil, err
	}

	return cartItems, nil
}

func (u *userService) CreateOrder(ctx context.Context, uId uint, orderRef string, pId string, amount float64) error {
	return nil
}

func (u *userService) CreateProfile(ctx context.Context, id uint, input *dto.Profile) error {
	user, err := u.findUserById(ctx, id)
	if err != nil {
		slg.Logger.Error("failed to find user", "id", id, "error", err)
		return err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}

	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = u.userRepository.UpdateUser(ctx, id, user)

	address := domain.Address{
		AddressLine1: input.Address.AddressLine1,
		AddressLine2: input.Address.AddressLine2,
		City:         input.Address.City,
		Country:      input.Address.Country,
		PostCode:     input.Address.PostCode,
		UserId:       id,
	}

	if err = u.userRepository.CreateProfile(ctx, &address); err != nil {
		return err
	}

	return nil
}

func (u *userService) GetProfile(ctx context.Context, id uint) (*domain.User, error) {
	profile, err := u.findUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (u *userService) UpdateProfile(ctx context.Context, id uint, input *dto.Profile) error {
	user, err := u.findUserById(ctx, id)
	if err != nil {
		return err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}

	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = u.userRepository.UpdateUser(ctx, id, user)
	if err != nil {
		return err
	}

	address := domain.Address{
		AddressLine1: input.Address.AddressLine1,
		AddressLine2: input.Address.AddressLine2,
		City:         input.Address.City,
		Country:      input.Address.Country,
		PostCode:     input.Address.PostCode,
		UserId:       id,
	}

	if err = u.userRepository.UpdateProfile(ctx, &address); err != nil {
		return err
	}

	return nil
}

func (u *userService) GetOrderById(ctx context.Context, id, uId uint) (*domain.Order, error) {
	return nil, nil
}

func (u *userService) GetOrders(ctx context.Context, user *domain.User) ([]*domain.Order, error) {
	return nil, nil
}

func NewUserService(userRepository repository.UserRepository, catalogRepository repository.CatalogRepository, authentication helper.Authentication, email notification.Email) UserService {
	return &userService{
		userRepository:    userRepository,
		catalogRepository: catalogRepository,
		authentication:    authentication,
		email:             email,
	}
}
