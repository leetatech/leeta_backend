package mailer

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	postmarkAPIURL                = "https://api.postmarkapp.com"
	sendEmailWithTemplateEndpoint = "/email/withTemplate/"
	fromLeeta                     = "admin@getleeta.com"
)

// Ensure implementation of EmailService interface
var _ Client = (*EmailService)(nil)

type EmailService struct {
	RESTClient *resty.Client
}

type Client interface {
	SendWithTemplate(message models.Message) error
}

func New(postmarkServerToken string) Client {
	// Build REST client
	restClient := resty.New()
	restClient.SetBaseURL(postmarkAPIURL)
	restClient.SetHeader("Content-Type", "application/json")
	restClient.SetHeader("Accept", "application/json")
	restClient.SetHeader("X-Postmark-Server-Token", postmarkServerToken)
	restClient.SetDebug(true)

	client := EmailService{
		RESTClient: restClient,
	}
	return &client
}

func (c *EmailService) SendWithTemplate(message models.Message) error {
	if c.RESTClient != nil {
		log.Debug().Msg("RESTClient is initialized")
	} else {
		log.Debug().Msg("RESTClient is not initialized")
	}

	if c.RESTClient == nil {
		return fmt.Errorf("RESTClient is not initialized")
	}

	payload := map[string]interface{}{
		"From":          fromLeeta,
		"To":            message.Target,
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
		return fmt.Errorf("failed to send email: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("resp is nil")
	}

	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			return fmt.Errorf("email template not found")
		}
		return fmt.Errorf("failed to send email: %s", resp.Status())
	}

	return nil
}
