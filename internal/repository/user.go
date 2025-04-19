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

func NewUserRepository(dbWrite, dbRead *gorm.DB) UserRepository {
	return &userRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
