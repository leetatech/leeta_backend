package messaging

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/leetError"
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
	SES     *ses.SES
	SNS     *sns.SNS
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
		log.Println("Error occurred while creating awsEmail session", err)
		return leetError.ErrorResponseBody(leetError.AwsSessionError, err)
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
