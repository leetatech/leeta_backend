package sms

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/pkg/notification"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/rs/zerolog/log"
)

type Client struct {
	Client notification.AWSClient
}

func New(awsClient notification.AWSClient) Client {
	return Client{Client: awsClient}
}

func (awsClient Client) SendSMS(message models.Message) error {
	input := &sns.PublishInput{
		PhoneNumber: aws.String(message.Target),
		Message:     aws.String(message.Body),
	}

	_, err := awsClient.Client.SNS.Publish(input) // TODO: enhance response validation
	if err != nil {
		log.Error().Err(err).Msg("Failed to send SMS")
		return errs.Body(errs.SnsSendSMSError, err)
	}

	log.Debug().Msg("SMS sent successfully")
	return nil
}
