package awsEmail

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/leetatech/leeta_backend/pkg/leetError"
	"github.com/leetatech/leeta_backend/pkg/messaging"
	"github.com/leetatech/leeta_backend/pkg/messaging/mailer/parseTemplates"
	"github.com/leetatech/leeta_backend/services/models"
	"go.uber.org/zap"
)

type AWSEmailClient struct {
	Client messaging.AWSClient
}

func NewAWSEmailClient(awsClient messaging.AWSClient) AWSEmailClient {
	return AWSEmailClient{Client: awsClient}
}

func (awsClient AWSEmailClient) SendEmail(templatePath string, message models.Message) error {
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
		Source: aws.String(fmt.Sprintf("Leeta Technologies <%s>", message.Sender)),
	}

	resp, err := awsClient.Client.SES.SendEmail(input)
	if err != nil {
		awsClient.Client.Log.Error("Error sending mail - ", zap.Error(err))
		return leetError.ErrorResponseBody(leetError.SesSendEmailError, err)
	}

	awsClient.Client.Log.Info("Email sent successfully", zap.String("MessageId", *resp.MessageId))

	return nil
}
