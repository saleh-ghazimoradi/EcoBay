package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserById(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, id uint, user *domain.User) (*domain.User, error)
}

type userRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := u.dbWrite.WithContext(ctx).Create(&user).Error; err != nil {
		slg.Logger.Error("create user", "error", err)
		return nil, ErrsCreate
	}
	return user, nil
}

func (u *userRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := u.dbRead.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		slg.Logger.Error("find user by email", "error", err)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrNotFound
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
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, id uint, user *domain.User) (*domain.User, error) {
	result := u.dbWrite.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Updates(user)

	if result.Error != nil {
		slg.Logger.Error("update user", "id", id, "error", result.Error)
		if errors.Is(result.Error, context.Canceled) {
			return nil, context.Canceled
		}
		return nil, ErrUpdate
	}

	if result.RowsAffected == 0 {
		slg.Logger.Warn("user not found for update", "id", id)
		return nil, ErrNotFound
	}

	var updatedUser domain.User
	if err := u.dbRead.WithContext(ctx).First(&updatedUser, id).Error; err != nil {
		slg.Logger.Error("failed to fetch updated user", "id", id, "error", err)
		return nil, ErrUpdate
	}
	return &updatedUser, nil
}

func NewUserRepository(dbWrite, dbRead *gorm.DB) UserRepository {
	return &userRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
