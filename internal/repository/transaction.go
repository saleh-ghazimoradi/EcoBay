package repository

import (
	"context"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	FindInitialPayment(ctx context.Context, uId uint) (*domain.Payment, error)
	FindOrders(ctx context.Context, uId uint) ([]*domain.OrderItem, error)
	FindOrderById(ctx context.Context, uId uint, id uint) (*dto.SellerOrderDetails, error)
	UpdatePayment(ctx context.Context, payment *domain.Payment) error
}

type transactionRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (t *transactionRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	return t.dbWrite.Create(payment).Error
}

func (t *transactionRepository) FindInitialPayment(ctx context.Context, uId uint) (*domain.Payment, error) {
	var payment *domain.Payment
	err := t.dbRead.WithContext(ctx).First(&payment, "user_id=? AND status=?", uId, "initial").Order("created_at desc").Error
	return payment, err
}

func (t *transactionRepository) FindOrders(ctx context.Context, uId uint) ([]*domain.OrderItem, error) {
	var orderItems []*domain.OrderItem
	if err := t.dbRead.WithContext(ctx).Find(&orderItems, "user_id = ?", uId).Error; err != nil {
		return nil, customErr.ErrNotFound
	}
	return orderItems, nil
}

func (t *transactionRepository) FindOrderById(ctx context.Context, uId uint, id uint) (*dto.SellerOrderDetails, error) {
	var sellerOrderItem dto.SellerOrderDetails
	if err := t.dbRead.WithContext(ctx).First(&sellerOrderItem, "user_id = ? AND id = ?", uId, id).Error; err != nil {
		return nil, customErr.ErrNotFound
	}
	return &sellerOrderItem, nil
}

func (t *transactionRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	return t.dbWrite.WithContext(ctx).Save(payment).Error
}

func NewTransactionRepository(dbWrite *gorm.DB, dbRead *gorm.DB) TransactionRepository {
	return &transactionRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
