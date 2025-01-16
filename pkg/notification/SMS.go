package notification

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
)

type NotificationsClient interface {
	SendSMS(ctx context.Context, phone string, message string) error
}

type notificationsClient struct {
}

func (n *notificationsClient) SendSMS(ctx context.Context, phone string, message string) error {
	accountSid := config.AppConfig.Necessities.AccountSMSSid
	authToken := config.AppConfig.Necessities.AuthToken

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(config.AppConfig.Necessities.SetFROM)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("failed to send SMS to %s: %v", phone, err)
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	log.Printf("SMS sent to %s: SID=%s, Status=%s", phone, *resp.Sid, *resp.Status)
	return nil
}

func NewNotificationsClient() NotificationsClient {
	return &notificationsClient{}
}
