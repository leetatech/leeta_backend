package mailer

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.uber.org/zap"
)

const (
	postmarkAPIURL                = "https://api.postmarkapp.com"
	sendEmailWithTemplateEndpoint = "/email/withTemplate/"
	fromLeeta                     = "admin@getleeta.com"
)

// Ensure implementation of EmailClient interface
var _ MailerClient = (*emailClient)(nil)

type emailClient struct {
	RESTClient *resty.Client
	logger     *zap.Logger
}

type MailerClient interface {
	SendEmailWithTemplate(message models.Message) error
}

func NewMailerClient(postmarkServerToken string, logger *zap.Logger) MailerClient {
	// Build REST client
	restClient := resty.New()
	restClient.SetBaseURL(postmarkAPIURL)
	restClient.SetHeader("Content-Type", "application/json")
	restClient.SetHeader("Accept", "application/json")
	restClient.SetHeader("X-Postmark-Server-Token", postmarkServerToken)
	restClient.SetDebug(true)

	emailClient := emailClient{
		RESTClient: restClient,
		logger:     logger,
	}

	return &emailClient
}

func (c *emailClient) SendEmailWithTemplate(message models.Message) error {
	payload := map[string]interface{}{
		"From":          fromLeeta,
		"To":            message.Target,
		"Subject":       message.Title,
		"TemplateAlias": message.TemplateID,
		"TemplateModel": message.DataMap,
	}

	var (
		result        EmailWithTemplateResponse
		errorResponse ErrorResponse
	)
	resp, err := c.RESTClient.R().
		SetBody(payload).
		SetResult(&result).
		SetError(&errorResponse).
		Post(sendEmailWithTemplateEndpoint)
	if err != nil {
		c.logger.Error("failed to send email", zap.Error(err))
		return fmt.Errorf("failed to send email: %w", err)
	}

	if resp.IsError() {
		c.logger.Error("failed to send email", zap.Any("error", errorResponse))
		return fmt.Errorf("failed to send email: %s", resp.Status())
	}

	c.logger.Info("email sent successfully", zap.Any("response", result))

	return nil
}
