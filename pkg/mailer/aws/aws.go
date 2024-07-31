package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/leetatech/leeta_backend/pkg/config"
	"github.com/leetatech/leeta_backend/pkg/mailer/parseTemplates"
	"github.com/leetatech/leeta_backend/services/models"
	"net/http"
	"time"
)

type MailClient struct {
	Config  *config.AWSConfig
	Session *session.Session
	SVC     *ses.SES
}

func (client *MailClient) Connect() error {
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	awsConfig := &aws.Config{
		Region: aws.String(client.Config.Region),
		Credentials: credentials.NewStaticCredentials(
			client.Config.Endpoint,
			client.Config.Secret,
			"",
		),
		HTTPClient: httpClient,
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	client.Session = sess
	client.SVC = ses.New(sess)

	return nil
}

func (client *MailClient) SendEmail(templatePath string, message models.Message) error {
	templateBody, err := parseTemplates.CreateSingleTemplate(templatePath, message)
	if err != nil {
		return fmt.Errorf("failed to create email template: %w", err)
	}

	if message.Sender == "" {
		return fmt.Errorf("sender address is empty")
	}

	if len(message.Recipients) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	validRecipients := filterValidEmails(message.Recipients)
	if len(validRecipients) == 0 {
		return fmt.Errorf("no valid recipients specified")
	}

	validCcRecipients := filterValidEmails(message.CcRecipients)
	validBccRecipients := filterValidEmails(message.BccRecipients)

	emailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			BccAddresses: toStringPointerSlice(validBccRecipients),
			CcAddresses:  toStringPointerSlice(validCcRecipients),
			ToAddresses:  toStringPointerSlice(validRecipients),
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

	_, err = client.SVC.SendEmail(emailInput)
	if err != nil {
		return fmt.Errorf("failed to send email using aws: %w", err)
	}

	return nil
}

func filterValidEmails(emails []string) []string {
	var validEmails []string
	for _, email := range emails {
		if email != "" {
			validEmails = append(validEmails, email)
		}
	}
	return validEmails
}

func toStringPointerSlice(strings []string) []*string {
	pointerSlice := make([]*string, len(strings))
	for i, s := range strings {
		pointerSlice[i] = aws.String(s)
	}
	return pointerSlice
}
