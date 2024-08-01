package notification

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/models"
	"log"
	"net/http"
	"time"
)

type AWSClient struct {
	Config  *config.AWSConfig
	Session *session.Session
	SES     *ses.SES
	SNS     *sns.SNS
}

func (awsClient *AWSClient) Connect() error {
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	awsConfig := &aws.Config{
		Region: aws.String(awsClient.Config.Region),
		Credentials: credentials.NewStaticCredentials(
			awsClient.Config.Endpoint,
			awsClient.Config.Secret,
			"",
		),
		HTTPClient: httpClient,
	}
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		log.Println("Error occurred while creating awsEmail session", err)
		return errs.Body(errs.AwsSessionError, err)
	}

	awsClient.Session = awsSession
	awsClient.SES = ses.New(awsClient.Session)
	awsClient.SNS = sns.New(awsClient.Session)

	return nil
}

type AWSClientInterface interface {
	SendEmail(templatePath string, message models.Message) error
	SendSMS(message models.Message) error
}
