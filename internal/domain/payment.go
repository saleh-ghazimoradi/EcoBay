package domain

import "time"

type PaymentStatus string

const (
	PaymentStatusInitial PaymentStatus = "initial"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
)

type Payment struct {
	ID            uint          `gorm:"PrimaryKey" json:"id"`
	UserId        uint          `json:"user_id"`
	CaptureMethod string        `json:"capture_method"`
	Amount        float64       `json:"amount"`
	OrderId       string        `json:"order_id"`
	CustomerId    string        `json:"customer_id"`
	PaymentId     string        `json:"payment_id"`
	Status        PaymentStatus `json:"status" gorm:"default:initial"`
	Response      string        `json:"response"`
	ClientSecret  string        `json:"client_secret"`
	CreatedAt     time.Time     `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"default:current_timestamp" json:"updated_at"`
}
