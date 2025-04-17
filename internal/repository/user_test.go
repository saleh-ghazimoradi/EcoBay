package repository

import (
	"context"
	"testing"
	"time"

	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupDB(t *testing.T) (*gorm.DB, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err, "should create in-memory database")

	err = db.AutoMigrate(&domain.User{})
	assert.NoError(t, err, "should migrate schema")

	return db, db
}

// cleanupDB drops the schema
func cleanupDB(t *testing.T, db *gorm.DB) {
	err := db.Migrator().DropTable(&domain.User{})
	assert.NoError(t, err, "should drop schema")
}

func TestUserRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *domain.User
		ctxSetup      func() context.Context
		expectedError error
		verifyUser    func(t *testing.T, user *domain.User)
	}{
		{
			name: "successful creation",
			user: &domain.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Phone:     "1234567890",
				Password:  "password123",
				UserType:  "seller",
			},
			ctxSetup: func() context.Context {
				return context.Background()
			},
			expectedError: nil,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.NotZero(t, user.ID, "should assign an ID")
				assert.Equal(t, "John", user.FirstName, "should save first name")
				assert.Equal(t, "Doe", user.LastName, "should save last name")
				assert.Equal(t, "john.doe@example.com", user.Email, "should save email")
				assert.Equal(t, "1234567890", user.Phone, "should save phone")
				assert.Equal(t, "password123", user.Password, "should save password")
				assert.Equal(t, "seller", user.UserType, "should save user type")
				assert.NotZero(t, user.CreatedAt, "should set created at")
				assert.NotZero(t, user.UpdatedAt, "should set updated at")
			},
		},
		{
			name: "duplicate email",
			user: &domain.User{
				Email:    "john.doe@example.com",
				Password: "password123",
			},
			ctxSetup: func() context.Context {
				return context.Background()
			},
			expectedError: ErrsCreate,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.Nil(t, user, "should return nil user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbWrite, dbRead := setupDB(t)
			defer cleanupDB(t, dbWrite)

			repo := NewUserRepository(dbWrite, dbRead)

			if tt.name == "duplicate email" {
				_, err := repo.CreateUser(context.Background(), &domain.User{
					Email:    "john.doe@example.com",
					Password: "password123",
				})
				assert.NoError(t, err, "should create first user")
			}

			ctx := tt.ctxSetup()
			result, err := repo.CreateUser(ctx, tt.user)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError, "should return correct error")
				tt.verifyUser(t, result)
			} else {
				assert.NoError(t, err, "should not return an error")
				assert.NotNil(t, result, "should return a user")
				tt.verifyUser(t, result)
			}
		})
	}
}

func TestUserRepository_FindUserByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		setup         func(t *testing.T, repo UserRepository)
		ctxSetup      func() context.Context
		expectedError error
		verifyUser    func(t *testing.T, user *domain.User)
	}{
		{
			name:  "find existing user",
			email: "test@example.com",
			setup: func(t *testing.T, repo UserRepository) {
				_, err := repo.CreateUser(context.Background(), &domain.User{
					Email:    "test@example.com",
					Password: "password123",
					UserType: "buyer",
				})
				assert.NoError(t, err, "should create user")
			},
			ctxSetup: func() context.Context {
				return context.Background()
			},
			expectedError: nil,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.NotNil(t, user, "should return a user")
				assert.Equal(t, "test@example.com", user.Email, "should return correct email")
				assert.Equal(t, "password123", user.Password, "should return correct password")
				assert.Equal(t, "buyer", user.UserType, "should return correct user type")
			},
		},
		{
			name:          "non-existing user",
			email:         "missing@example.com",
			setup:         func(t *testing.T, repo UserRepository) {},
			ctxSetup:      func() context.Context { return context.Background() },
			expectedError: ErrNotFound,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.Nil(t, user, "should return nil user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbWrite, dbRead := setupDB(t)
			defer cleanupDB(t, dbWrite)

			repo := NewUserRepository(dbWrite, dbRead)

			tt.setup(t, repo)

			ctx := tt.ctxSetup()
			user, err := repo.FindUserByEmail(ctx, tt.email)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError, "should return correct error")
				tt.verifyUser(t, user)
			} else {
				assert.NoError(t, err, "should not return an error")
				tt.verifyUser(t, user)
			}
		})
	}
}

func TestUserRepository_FindUserById(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		setup         func(t *testing.T, repo UserRepository)
		ctxSetup      func() context.Context
		expectedError error
		verifyUser    func(t *testing.T, user *domain.User)
	}{
		{
			name: "find existing user",
			id:   1,
			setup: func(t *testing.T, repo UserRepository) {
				_, err := repo.CreateUser(context.Background(), &domain.User{
					Email:    "test@example.com",
					Password: "password123",
					UserType: "buyer",
				})
				assert.NoError(t, err, "should create user")
			},
			ctxSetup: func() context.Context {
				return context.Background()
			},
			expectedError: nil,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.NotNil(t, user, "should return a user")
				assert.Equal(t, uint(1), user.ID, "should return correct ID")
				assert.Equal(t, "test@example.com", user.Email, "should return correct email")
				assert.Equal(t, "password123", user.Password, "should return correct password")
				assert.Equal(t, "buyer", user.UserType, "should return correct user type")
			},
		},
		{
			name:          "non-existing user",
			id:            999,
			setup:         func(t *testing.T, repo UserRepository) {},
			ctxSetup:      func() context.Context { return context.Background() },
			expectedError: ErrNotFound,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.Nil(t, user, "should return nil user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbWrite, dbRead := setupDB(t)
			defer cleanupDB(t, dbWrite)

			repo := NewUserRepository(dbWrite, dbRead)

			tt.setup(t, repo)

			ctx := tt.ctxSetup()
			user, err := repo.FindUserById(ctx, tt.id)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError, "should return correct error")
				tt.verifyUser(t, user)
			} else {
				assert.NoError(t, err, "should not return an error")
				tt.verifyUser(t, user)
			}
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		user          *domain.User
		setup         func(t *testing.T, repo UserRepository)
		ctxSetup      func() context.Context
		expectedError error
		verifyUser    func(t *testing.T, user *domain.User)
	}{
		{
			name: "update all fields",
			id:   1,
			user: &domain.User{
				FirstName: "Jane",
				LastName:  "Smith",
				Phone:     "9876543210",
				UserType:  "seller",
				Email:     "jane.smith@example.com",
				Password:  "newpassword",
				Code:      654321,
				Expiry:    time.Now().Add(48 * time.Hour),
				Verified:  true,
			},
			setup: func(t *testing.T, repo UserRepository) {
				_, err := repo.CreateUser(context.Background(), &domain.User{
					Email:    "test@example.com",
					Password: "password123",
					UserType: "buyer",
				})
				assert.NoError(t, err, "should create user")
			},
			ctxSetup: func() context.Context {
				return context.Background()
			},
			expectedError: nil,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.NotNil(t, user, "should return a user")
				assert.Equal(t, uint(1), user.ID, "should keep same ID")
				assert.Equal(t, "Jane", user.FirstName, "should update first name")
				assert.Equal(t, "Smith", user.LastName, "should update last name")
				assert.Equal(t, "9876543210", user.Phone, "should update phone")
				assert.Equal(t, "seller", user.UserType, "should update user type")
				assert.Equal(t, "jane.smith@example.com", user.Email, "should update email")
				assert.Equal(t, "newpassword", user.Password, "should update password")
				assert.Equal(t, 654321, user.Code, "should update code")
				assert.WithinDuration(t, time.Now().Add(48*time.Hour), user.Expiry, time.Minute, "should update expiry")
				assert.True(t, user.Verified, "should update verified")
				assert.NotZero(t, user.UpdatedAt, "should update updated at")
			},
		},
		{
			name: "update partial fields",
			id:   1,
			user: &domain.User{
				FirstName: "Jane",
				Phone:     "9876543210",
			},
			setup: func(t *testing.T, repo UserRepository) {
				_, err := repo.CreateUser(context.Background(), &domain.User{
					Email:    "test@example.com",
					Password: "password123",
					UserType: "buyer",
					LastName: "Doe",
				})
				assert.NoError(t, err, "should create user")
			},
			ctxSetup: func() context.Context {
				return context.Background()
			},
			expectedError: nil,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.NotNil(t, user, "should return a user")
				assert.Equal(t, uint(1), user.ID, "should keep same ID")
				assert.Equal(t, "Jane", user.FirstName, "should update first name")
				assert.Equal(t, "Doe", user.LastName, "should keep last name")
				assert.Equal(t, "9876543210", user.Phone, "should update phone")
				assert.Equal(t, "buyer", user.UserType, "should keep user type")
				assert.Equal(t, "test@example.com", user.Email, "should keep email")
				assert.Equal(t, "password123", user.Password, "should keep password")
				assert.Zero(t, user.Code, "should keep code")
				assert.True(t, user.Expiry.IsZero(), "should keep expiry")
				assert.False(t, user.Verified, "should keep verified")
				assert.NotZero(t, user.UpdatedAt, "should update updated at")
			},
		},
		{
			name: "non-existing user",
			id:   999,
			user: &domain.User{
				FirstName: "Jane",
				Phone:     "9876543210",
			},
			setup:         func(t *testing.T, repo UserRepository) {},
			ctxSetup:      func() context.Context { return context.Background() },
			expectedError: ErrNotFound,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.Nil(t, user, "should return nil user")
			},
		},
		{
			name: "context canceled",
			id:   1,
			user: &domain.User{
				FirstName: "Jane",
				Phone:     "9876543210",
			},
			setup: func(t *testing.T, repo UserRepository) {
				_, err := repo.CreateUser(context.Background(), &domain.User{
					Email:    "test@example.com",
					Password: "password123",
				})
				assert.NoError(t, err, "should create user")
			},
			ctxSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			expectedError: context.Canceled,
			verifyUser: func(t *testing.T, user *domain.User) {
				assert.Nil(t, user, "should return nil user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbWrite, dbRead := setupDB(t)
			defer cleanupDB(t, dbWrite)

			repo := NewUserRepository(dbWrite, dbRead)

			tt.setup(t, repo)

			ctx := tt.ctxSetup()
			user, err := repo.UpdateUser(ctx, tt.id, tt.user)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError, "should return correct error")
				tt.verifyUser(t, user)
			} else {
				assert.NoError(t, err, "should not return an error")
				tt.verifyUser(t, user)
			}
		})
	}
}
