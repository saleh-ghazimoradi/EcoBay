package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserById(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, id uint, user *domain.User) (*domain.User, error)
	CreateBankAccount(ctx context.Context, bankAccount *domain.BankAccount) error
	FindCartItems(ctx context.Context, uId uint) ([]*domain.Cart, error)
	FindCartItem(ctx context.Context, uId uint, pId uint) (*domain.Cart, error)
	CreateCart(ctx context.Context, cart *domain.Cart) error
	UpdateCart(ctx context.Context, cart *domain.Cart) error
	DeleteCartById(ctx context.Context, id uint) error
	DeleteCartItems(ctx context.Context, uId uint) error
	CreateOrder(ctx context.Context, order *domain.Order) error
	FindOrders(ctx context.Context, uId uint) ([]*domain.Order, error)
	FindOrderById(ctx context.Context, id, uId uint) (*domain.Order, error)
	CreateProfile(ctx context.Context, address *domain.Address) error
	UpdateProfile(ctx context.Context, address *domain.Address) error
}

type userRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := u.dbWrite.WithContext(ctx).Create(&user).Error; err != nil {
		slg.Logger.Error("create user", "error", err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, customErr.ErrEmailExists
		}
		return nil, err
	}
	return user, nil
}

func (u *userRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := u.dbRead.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		slg.Logger.Error("find user by email", "error", err)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, customErr.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *userRepository) FindUserById(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := u.dbRead.WithContext(ctx).First(&user, id).Error; err != nil {
		slg.Logger.Error("find user by id", "error", err)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, customErr.ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, id uint, user *domain.User) (*domain.User, error) {
	if err := u.dbWrite.WithContext(ctx).Model(&user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(user).Error; err != nil {
		slg.Logger.Error("update user", "error", err)
		return nil, customErr.ErrUpdate
	}
	return user, nil
}

func (u *userRepository) CreateBankAccount(ctx context.Context, bankAccount *domain.BankAccount) error {
	if err := u.dbWrite.WithContext(ctx).Create(&bankAccount).Error; err != nil {
		slg.Logger.Error("create bank account", "error", err)
		return err
	}
	return nil
}

func (u *userRepository) FindCartItems(ctx context.Context, uId uint) ([]*domain.Cart, error) {
	var carts []*domain.Cart
	err := u.dbRead.WithContext(ctx).Where("user_id = ?", uId).Find(&carts).Error
	return carts, err
}

func (u *userRepository) FindCartItem(ctx context.Context, uId uint, pId uint) (*domain.Cart, error) {
	var cartItem domain.Cart
	err := u.dbRead.WithContext(ctx).Where("user_id=? AND product_id=?", uId, pId).First(&cartItem).Error
	return &cartItem, err
}

func (u *userRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	if err := u.dbWrite.WithContext(ctx).Create(&cart).Error; err != nil {
		slg.Logger.Error("create cart item", "error", err)
		return customErr.ErrsCreate
	}
	return nil
}

func (u *userRepository) UpdateCart(ctx context.Context, cart *domain.Cart) error {
	var c domain.Cart
	if err := u.dbWrite.WithContext(ctx).Model(&c).Clauses(clause.Returning{}).Where("id=?", cart.ID).Updates(cart).Error; err != nil {
		return customErr.ErrUpdate
	}
	return nil
}

func (u *userRepository) DeleteCartById(ctx context.Context, id uint) error {
	if err := u.dbWrite.WithContext(ctx).Delete(&domain.Cart{}, id).Error; err != nil {
		return customErr.ErrDelete
	}
	return nil
}

func (u *userRepository) DeleteCartItems(ctx context.Context, uId uint) error {
	if err := u.dbWrite.WithContext(ctx).Where("user_id = ?", uId).Delete(&domain.Cart{}).Error; err != nil {
		slg.Logger.Error("delete cart items", "error", err)
		return err
	}
	return nil
}

func (u *userRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	return nil
}

func (u *userRepository) FindOrders(ctx context.Context, uId uint) ([]*domain.Order, error) {
	return nil, nil
}

func (u *userRepository) FindOrderById(ctx context.Context, id, uId uint) (*domain.Order, error) {
	return nil, nil
}

func (u *userRepository) CreateProfile(ctx context.Context, address *domain.Address) error {
	return nil
}

func (u *userRepository) UpdateProfile(ctx context.Context, address *domain.Address) error {
	return nil
}

func NewUserRepository(dbWrite, dbRead *gorm.DB) UserRepository {
	return &userRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
