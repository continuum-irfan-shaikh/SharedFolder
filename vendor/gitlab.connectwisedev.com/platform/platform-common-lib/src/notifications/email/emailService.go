package email

import (
	"errors"
	"fmt"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"strings"
)

const (
	DefaultRegion  string         = "us-east-1"
	DefaultCharset string         = "UTF-8"
	FULL_NAME      PlaceholderKey = "{fullName}"
	EMAIL_BODY     PlaceholderKey = "{emailBody}"
	ButtonLabel    PlaceholderKey = "{buttonLabel}"
	LINK           PlaceholderKey = "{link}"
)

func (ec *EmailContent) AddRecipientEmail(recipeint string) *EmailContent {
	if recipeint != "" {
		ec.ToAddresses = append(ec.ToAddresses, &recipeint)
	}
	return ec
}

func (ec *EmailContent) AddCCEmail(ccaddr string) *EmailContent {
	if ccaddr != "" {
		ec.CCAddresses = append(ec.CCAddresses, &ccaddr)
	}
	return ec
}
func (ec *EmailContent) AddTemplateKey(key PlaceholderKey, value string) *EmailContent {
	if ec.ContentKeyValue == nil {
		ec.ContentKeyValue = make(map[PlaceholderKey]string)
	}
	if key != "" {
		ec.ContentKeyValue[key] = value
	}

	return ec
}
func (ec *EmailContent) AddBodyKey(key PlaceholderKey, value string) *EmailContent {
	if ec.BodyKeyValue == nil {
		ec.BodyKeyValue = make(map[PlaceholderKey]string)
	}
	if key != "" {
		ec.BodyKeyValue[key] = value
	}
	return ec
}

func (ec *EmailContent) replaceTemplate() {
	if len(ec.HTMLTemplate) > 0 {
		for k, v := range ec.ContentKeyValue {
			ec.HTMLTemplate = strings.ReplaceAll(ec.HTMLTemplate, string(k), v)
		}
	}
}

func (ec *EmailContent) replaceBody() {
	if len(ec.HTMLBody) > 0 {
		for k, v := range ec.BodyKeyValue {
			ec.HTMLBody = strings.ReplaceAll(ec.HTMLBody, string(k), v)
		}
		ec.HTMLTemplate = strings.ReplaceAll(ec.HTMLTemplate, string(EMAIL_BODY), ec.HTMLBody)
	}
}

func (ec *EmailContent) validate() error {
	if len(ec.ToAddresses) <= 0 {
		return errors.New("Please provide valid recepient(s)")
	}
	if len(ec.Sender) == 0 {
		return errors.New("Please provide valid source")
	}
	return nil
}

//Composition of SES Client
type AWSEmailClient struct {
	Region    string
	SESClient sesiface.SESAPI
	Log       logger.Log
}

//Configure AwS SES Client
func (ases *AWSEmailClient) Configure() error {
	if ases.Region == "" {
		ases.Region = DefaultRegion
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ases.Region)},
	)
	if err != nil {
		return err
	}
	// Create an SES session.
	svc := ses.New(sess)
	ases.SESClient = svc
	return nil
}

func (ases *AWSEmailClient) buildMessage(ec *EmailContent) *ses.SendEmailInput {
	ec.replaceBody()
	ec.replaceTemplate()
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: ec.ToAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(ec.CharSet),
					Data:    aws.String(ec.HTMLTemplate),
				},
				Text: &ses.Content{
					Charset: aws.String(ec.CharSet),
					Data:    aws.String(ec.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(ec.CharSet),
				Data:    aws.String(ec.Subject),
			},
		},
		Source: aws.String(ec.Sender),
	}
	return input
}

//Error Handler passed to handle AWS error codes.
func (ases *AWSEmailClient) AWSErrorHandler(transactionID string, err error) {
	if aerr, ok := err.(awserr.Error); ok {
		ases.logError(transactionID, aerr.Code(), aerr.Error())
	} else {
		ases.logError(transactionID, ses.ErrCodeMessageRejected, aerr.Error())
	}
}

func (ases *AWSEmailClient) logError(transactionID string, errorCode string, errorMessage string) {
	if ases.Log != nil {
		ases.Log.Info(transactionID, errorCode, errorMessage)
	} else {
		fmt.Println(transactionID, errorCode, errorMessage)
	}
}

func (ases *AWSEmailClient) SendEmail(transactionID string, ec *EmailContent, errorHandlerCallback ErrorHandler) (*EmailOutput, error) {
	// Attempt to send the email.
	err := ec.validate()
	if err != nil {
		return nil, err
	}
	result, err := ases.SESClient.SendEmail(ases.buildMessage(ec))
	if err != nil {
		if errorHandlerCallback != nil {
			errorHandlerCallback(transactionID, err)
		}
		return nil, err
	}

	return &EmailOutput{result.MessageId}, nil
}
