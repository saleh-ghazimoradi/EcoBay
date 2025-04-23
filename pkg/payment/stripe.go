package payment

import (
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

type PaymentsClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error)
	GetPaymentStatus(pId string) (*stripe.PaymentIntent, error)
}

type paymentsClient struct {
	stripeSecretKey string
	successUrl      string
	cancelUrl       string
}

func (p *paymentsClient) CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.stripeSecretKey
	amountInCents := amount * 100

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(amountInCents)),
		Currency:           stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	params.AddMetadata("order_id", fmt.Sprintf("%s", orderId))
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))

	pi, err := paymentintent.New(params)
	if err != nil {
		slg.Logger.Error("Error creating session", "error", err)
		return nil, errors.New("payment intent creation failed")
	}

	return pi, nil
}

func (p *paymentsClient) GetPaymentStatus(pId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.stripeSecretKey
	params := &stripe.PaymentIntentParams{}
	result, err := paymentintent.Get(pId, params)
	if err != nil {
		slg.Logger.Error("Error getting session", "error", err)
		return nil, errors.New("get payment intent failed")
	}
	return result, nil
}

func NewPaymentsClient(stripeSecretKey, successUrl, cancelUrl string) PaymentsClient {
	return &paymentsClient{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		cancelUrl:       cancelUrl,
	}
}
