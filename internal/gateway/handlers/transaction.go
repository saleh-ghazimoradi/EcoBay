package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"github.com/saleh-ghazimoradi/EcoBay/pkg/payment"
	"net/http"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	authentication     helper.Authentication
	paymentClient      payment.PaymentsClient
	userService        service.UserService
}

func (t *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {
	user := t.authentication.GetCurrentUser(ctx)
	pubKey := config.AppConfig.Stripe.PublishableKey

	activePayment, err := t.transactionService.GetActivePayment(context.Background(), user.ID)
	if err != nil {
		return serverErrorResponse(ctx, err)
	}

	if activePayment != nil {
		return successResponse(ctx, fiber.StatusOK, "payment successfully created", pubKey)
	}

	// No active payment exists, create a new one
	_, amount, err := t.userService.FindCart(context.Background(), user.ID)
	if err != nil {
		return serverErrorResponse(ctx, err)
	}

	orderId, err := helper.RandomNumbers(8)
	if err != nil {
		return serverErrorResponse(ctx, errors.New("error generating order id"))
	}

	paymentResult, err := t.paymentClient.CreatePayment(amount, user.ID, orderId)
	if err != nil {
		return badRequestResponse(ctx, err)
	}

	if err = t.transactionService.StoreCreatedPayment(context.Background(), &dto.CreatePayment{
		UserId:       user.ID,
		Amount:       amount,
		ClientSecret: paymentResult.ClientSecret,
		PaymentId:    paymentResult.ID,
		OrderId:      orderId,
	}); err != nil {
		return badRequestResponse(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "create payment",
		"pubKey":  pubKey,
		"secret":  paymentResult.ClientSecret,
	})
}

//func (t *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {
//	user := t.authentication.GetCurrentUser(ctx)
//
//	pubKey := config.AppConfig.Stripe.PublishableKey
//
//	activePayment, err := t.transactionService.GetActivePayment(context.Background(), user.ID)
//	if err != nil {
//		return serverErrorResponse(ctx, err)
//	}
//
//	if activePayment.ID > 0 {
//		return successResponse(ctx, fiber.StatusOK, "payment successfully created", pubKey)
//	}
//
//	_, amount, err := t.userService.FindCart(context.Background(), user.ID)
//
//	orderId, err := helper.RandomNumbers(8)
//	if err != nil {
//		return serverErrorResponse(ctx, errors.New("error generating order id"))
//	}
//
//	paymentResult, err := t.paymentClient.CreatePayment(amount, user.ID, orderId)
//	if err != nil {
//		return badRequestResponse(ctx, err)
//	}
//
//	if err = t.transactionService.StoreCreatedPayment(context.Background(), &dto.CreatePayment{
//		UserId:       user.ID,
//		Amount:       amount,
//		ClientSecret: paymentResult.ClientSecret,
//		PaymentId:    paymentResult.ID,
//		OrderId:      orderId,
//	}); err != nil {
//		return badRequestResponse(ctx, err)
//	}
//
//	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
//		"message": "create payment",
//		"pubKey":  pubKey,
//		"secret":  paymentResult.ClientSecret,
//	})
//}

func (t *TransactionHandler) VerifyPayment(ctx *fiber.Ctx) error {
	user := t.authentication.GetCurrentUser(ctx)

	activePayment, err := t.transactionService.GetActivePayment(context.Background(), user.ID)
	if err != nil || activePayment.ID == 0 {
		return notFoundResponse(ctx)
	}

	paymentRes, err := t.paymentClient.GetPaymentStatus(activePayment.PaymentId)
	if err != nil {
		return badRequestResponse(ctx, err)
	}

	paymentJson, _ := json.Marshal(paymentRes)
	paymentLogs := string(paymentJson)
	paymentStatus := "failed"

	if paymentRes.Status == "succeeded" {
		paymentStatus = "success"
		err = t.userService.CreateOrder(ctx.Context(), user.ID, activePayment.OrderId, activePayment.PaymentId, activePayment.Amount)
	}

	if err != nil {
		return serverErrorResponse(ctx, err)
	}

	if err = t.transactionService.UpdatePayment(ctx.Context(), user.ID, paymentStatus, paymentLogs); err != nil {
		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusOK, "payment successfully created", paymentRes)
}

func NewTransactionHandler(transactionService service.TransactionService, authentication helper.Authentication, paymentClient payment.PaymentsClient, userService service.UserService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		authentication:     authentication,
		paymentClient:      paymentClient,
		userService:        userService,
	}
}
