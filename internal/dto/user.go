package dto

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserSignup struct {
	UserLogin
	Phone string `json:"phone" validate:"required"`
}

type VerificationCode struct {
	Code string `json:"code" validate:"required,len=6"`
}
