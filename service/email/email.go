package email

import (
	"os"

	"github.com/resend/resend-go/v2"
)

type Email struct {
	ApiKey string
	From   string
	Client *resend.Client
}

func NewEmailService() *Email {
	return &Email{
		ApiKey: os.Getenv("RESEND_API_KEY"),
		From:   os.Getenv("EMAIL_FROM"),
		Client: resend.NewClient(os.Getenv("RESEND_API_KEY")),
	}
}

func (e Email) Send(to, subject, body string) error {

	params := &resend.SendEmailRequest{
		From:    e.From,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	_, err := e.Client.Emails.Send(params)
	return err
}

func (e Email) SendToMultipleRecipients(to []string, subject, body string) error {
	params := &resend.SendEmailRequest{
		From:    e.From,
		To:      to,
		Subject: subject,
		Html:    body,
	}
	_, err := e.Client.Emails.Send(params)
	return err
}
