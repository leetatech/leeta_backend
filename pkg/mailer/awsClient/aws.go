package awsClient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/mailer/parseTemplates"
	"github.com/leetatech/leeta_backend/services/models"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type AWSClient struct {
	Config  *config.ServerConfig
	Log     *zap.Logger
	Session *session.Session
	SVC     *ses.SES
}

func (awsClient *AWSClient) ConnectAWS() error {
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	awsConfig := &aws.Config{
		Region: aws.String(awsClient.Config.AWSConfig.Region),
		Credentials: credentials.NewStaticCredentials(
			awsClient.Config.AWSConfig.Endpoint,
			awsClient.Config.AWSConfig.Secret,
			"",
		),
		HTTPClient: httpClient,
	}
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		log.Println("Error occurred while creating aws session", err)
		return err
	}

	awsClient.Session = awsSession
	awsClient.SVC = ses.New(awsClient.Session)

	return nil
}

func (awsClient AWSClient) SendEmail(templatePath string, message models.Message) error {
	templateBody, err := parseTemplates.CreateSingleTemplate(templatePath, message)
	if err != nil {
		return err
	}
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			BccAddresses: message.BccRecipients,
			CcAddresses:  message.CcRecipients,
			ToAddresses:  message.Recipients,
		},

		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(templateBody),
				},
			},

			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(message.Title),
			},
		},

		Source: aws.String(message.Sender),
	}

	_, err = awsClient.SVC.SendEmail(input)
	if err != nil {
		log.Println("Error sending mail - ", err)
		return err
	}

	log.Println("Email sent successfully")

	return nil
}
