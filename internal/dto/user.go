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

type BecomeSellerInput struct {
	FirstName         string `json:"first_name" validate:"required,min=1,max=255"`
	LastName          string `json:"last_name" validate:"required,min=1,max=255"`
	PhoneNumber       string `json:"phone_number" validate:"required"`
	BankAccountNumber uint   `json:"bank_account_number" validate:"required"`
	SwiftCode         string `json:"swift_code" validate:"required"`
	PaymentType       string `json:"payment_type" validate:"required"`
}
