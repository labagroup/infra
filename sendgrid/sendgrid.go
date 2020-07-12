package sendgrid

import (
	"context"
	"fmt"

	"github.com/gopub/errors"
	"github.com/gopub/log"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Client struct {
	config *Config
	client *sendgrid.Client
}

type Config struct {
	APIKey      string
	SenderEmail string
	SenderName  string
}

func NewSendGridService(config *Config) *Client {
	if len(config.APIKey) == 0 {
		log.Fatal("Missing APIKey")
	}

	if len(config.SenderEmail) == 0 {
		log.Fatal("Missing SenderEmail")
	}

	if len(config.SenderName) == 0 {
		log.Fatal("Missing SenderName")
	}

	s := &Client{
		config: config,
		client: sendgrid.NewSendClient(config.APIKey),
	}
	return s
}

func (s *Client) Push(ctx context.Context, recipient, subject, content string) error {
	from := mail.NewEmail(s.config.SenderName, s.config.SenderEmail)
	to := mail.NewEmail(recipient, recipient)

	// TODO:
	// If email client doesn't support html content, plain text will be displayed,
	// so it'd be better pass plain text along with html
	plainTextContent := ""
	htmlContent := "<strong>and easy to do anywhere, even with Go!!!</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	resp, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	if resp.StatusCode >= 400 {
		return errors.Format(resp.StatusCode, resp.Body)
	}
	return nil
}

func (s *Client) PushTemplate(ctx context.Context, recipient, templateID string, params map[string]string) error {
	logger := log.FromContext(ctx).With("recipient", recipient, "template_id", templateID)
	// TODO:
	logger.Info("OK")
	return nil
}
