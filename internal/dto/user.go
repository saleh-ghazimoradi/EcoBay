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

type Address struct {
	AddressLine1 string `json:"address_line1" validate:"required,min=1,max=255"`
	AddressLine2 string `json:"address_line2" validate:"required,min=1,max=255"`
	City         string `json:"city" validate:"required,min=1,max=255"`
	PostCode     uint   `json:"post_code" validate:"required"`
	Country      string `json:"country" validate:"required"`
}

type Profile struct {
	FirstName string  `json:"first_name" validate:"required,min=1,max=255"`
	LastName  string  `json:"last_name" validate:"required,min=1,max=255"`
	Address   Address `json:"address" validate:"required"`
}
