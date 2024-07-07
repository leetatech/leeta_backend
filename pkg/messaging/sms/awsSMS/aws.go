package awsSMS

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/messaging"
	"github.com/leetatech/leeta_backend/services/models"
	"go.uber.org/zap"
)

type AWSSMSClient struct {
	Client messaging.AWSClient
}

func NewAWSSMSClient(awsClient messaging.AWSClient) AWSSMSClient {
	return AWSSMSClient{Client: awsClient}
}

func (awsClient AWSSMSClient) SendSMS(message models.Message) error {
	input := &sns.PublishInput{
		PhoneNumber: aws.String(message.Target),
		Message:     aws.String(message.Body),
	}

	resp, err := awsClient.Client.SNS.Publish(input)
	if err != nil {
		awsClient.Client.Log.Error("Error sending sms - ", zap.Error(err))
		return leetError.ErrorResponseBody(leetError.SnsSendSMSError, err)
	}

	awsClient.Client.Log.Info("sms sent successfully", zap.String("MessageId", *resp.MessageId))

	return nil
}
