package helper

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service/service_models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	Secret string
}

func (a *Auth) CreateHashedPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password length must be at least 6 characters long")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("password length must be at least 6 characters long")
	}

	return string(hashPassword), nil
}

func (a *Auth) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("required inputs are missing to generate token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", errors.New("error signing token")
	}

	return tokenStr, nil
}

func (a *Auth) VerifyPassword(password, hashedPassword string) error {
	if len(password) < 6 {
		return errors.New("password length must be at least 6 characters long")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return errors.New("password does not match")
	}

	return nil
}

func (a *Auth) VerifyToken(token string) (*service_models.User, error) {
	tokenArr := strings.Split(token, " ")
	if len(tokenArr) != 2 {
		return nil, errors.New("invalid token")
	}

	tokenStr := tokenArr[1]

	if tokenArr[0] != "Bearer" {
		return nil, errors.New("invalid token")
	}

	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.Secret), nil
	})

	if err != nil {
		return nil, errors.New("invalid signing method")
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("token is expired")
		}
		user := &service_models.User{}
		user.ID = uint(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.UserType = claims["role"].(string)
		return user, nil
	}
	return nil, errors.New("token verification failed")
}

func (a *Auth) Authorize(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	user, err := a.VerifyToken(authHeader)
	if err == nil && user.ID > 0 {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err,
		})
	}
}

func (a *Auth) GetCurrentUser(ctx *fiber.Ctx) (*service_models.User, error) {
	user := ctx.Locals("user")
	return user.(*service_models.User), nil
}

func NewAuth(secret string) Auth {
	return Auth{
		Secret: secret,
	}
}
