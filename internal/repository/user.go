package repository

import "gorm.io/gorm"

type UserRepository interface {
}

type userRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func NewUserRepository(dbWrite, dbRead *gorm.DB) UserRepository {
	return &userRepository{}
}
