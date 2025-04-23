package service

import (
	"context"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
)

type TransactionService interface {
	GetOrders(ctx context.Context, user *domain.User) ([]*domain.OrderItem, error)
	GetOrderDetails(ctx context.Context, user *domain.User, id uint) (*dto.SellerOrderDetails, error)
	GetActivePayment(ctx context.Context, uId uint) (*domain.Payment, error)
	StoreCreatedPayment(ctx context.Context, input *dto.CreatePayment) error
	UpdatePayment(ctx context.Context, uId uint, status, paymentLog string) error
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	authentication        helper.Authentication
}

func (t *transactionService) GetOrders(ctx context.Context, user *domain.User) ([]*domain.OrderItem, error) {
	orders, err := t.transactionRepository.FindOrders(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (t *transactionService) GetOrderDetails(ctx context.Context, user *domain.User, id uint) (*dto.SellerOrderDetails, error) {
	order, err := t.transactionRepository.FindOrderById(ctx, user.ID, id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (t *transactionService) GetActivePayment(ctx context.Context, uId uint) (*domain.Payment, error) {
	return t.transactionRepository.FindInitialPayment(ctx, uId)
}

func (t *transactionService) StoreCreatedPayment(ctx context.Context, input *dto.CreatePayment) error {
	payment := domain.Payment{
		UserId:       input.UserId,
		Amount:       input.Amount,
		Status:       domain.PaymentStatusInitial,
		PaymentId:    input.PaymentId,
		ClientSecret: input.ClientSecret,
		OrderId:      input.OrderId,
	}

	return t.transactionRepository.CreatePayment(ctx, &payment)
}

func (t *transactionService) UpdatePayment(ctx context.Context, uId uint, status, paymentLog string) error {
	p, err := t.GetActivePayment(ctx, uId)
	if err != nil {
		return err
	}

	p.Status = domain.PaymentStatus(status)
	p.Response = paymentLog

	return t.transactionRepository.UpdatePayment(ctx, p)
}

func NewTransactionService(transactionRepository repository.TransactionRepository, authentication helper.Authentication) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		authentication:        authentication,
	}
}
