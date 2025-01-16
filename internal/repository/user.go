package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service/service_models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *service_models.User) (*service_models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*service_models.User, error)
	FindUserByID(ctx context.Context, id uint) (*service_models.User, error)
	UpdateUser(ctx context.Context, id uint, user *service_models.User) (*service_models.User, error)
	CreateBankAccount(ctx context.Context, bankAccount *service_models.BankAccount) error
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) CreateUser(ctx context.Context, user *service_models.User) (*service_models.User, error) {
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, errors.New("create user failed")
	}
	return user, nil
}

func (u *userRepository) FindUserByEmail(ctx context.Context, email string) (*service_models.User, error) {
	var user service_models.User
	if err := u.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, errors.New("user does not exist")
	}

	return &user, nil
}

func (u *userRepository) FindUserByID(ctx context.Context, id uint) (*service_models.User, error) {
	var user service_models.User
	if err := u.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, errors.New("user does not exist")
	}

	return &user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, id uint, user *service_models.User) (*service_models.User, error) {
	if err := u.db.WithContext(ctx).Model(&user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return nil, errors.New("error updating user")
	}
	return user, nil
}

func (u *userRepository) CreateBankAccount(ctx context.Context, bankAccount *service_models.BankAccount) error {
	if err := u.db.WithContext(ctx).Create(&bankAccount).Error; err != nil {
		return errors.New("create bank account failed")
	}
	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
