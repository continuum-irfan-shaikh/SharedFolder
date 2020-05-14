package email

import (
	"errors"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/notifications/email/mocks"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	sender = "nitin.kothari@continuum.net"

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	recipient = "Recipient@continuum.net"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	subject = "Important Notification from Continuum"

	htmlTemplate = "<p>Welcome to ITSupport247. Your account has been successfully created. {emailBody} </p>"
	// The HTML body for the email.
	htmlBody = "<p>Login to the ITSupport247 Portal using the following information:</p>" +
		"<p>URL: {link}</p><p>User Name: {fullName}</p>" +
		"<p>Password: Click <a href=&quot;{link}&quot;>here</a> to create a new password for your account.</p>"

	//The email body for recipients with non-HTML email clients.
	textBody = "Welcome to ITS Portal. To log in to your ITSupport247 portal account, you will need to create a new password. Click on the link below to open a secure browser window and reset your password now."
)

func TestAWSEmailClient(t *testing.T) {
	content := &EmailContent{}
	content.AddRecipientEmail(recipient).AddRecipientEmail("nitin.kothari@continuum.net")
	content.CharSet = DefaultCharset
	content.TextBody = textBody
	content.Sender = sender
	content.Subject = subject
	content.HTMLBody = htmlBody
	content.HTMLTemplate = htmlTemplate
	content.AddCCEmail("dummy@dummy.com")
	content.AddTemplateKey(EMAIL_BODY, htmlBody)
	content.AddBodyKey(FULL_NAME, "nkothari")

	mockSESClient := &mocks.SESAPI{}

	awsClient := &AWSEmailClient{"us-east-1", nil, nil}
	awsClient.Configure()
	awsClient.SESClient = mockSESClient
	dummyId := "DummyID-1111"
	mockSESClient.On("SendEmail", mock.Anything).Return(&ses.SendEmailOutput{MessageId: &dummyId}, nil)
	emailOutput, _ := awsClient.SendEmail("", content, nil)
	assert.Equal(t, &dummyId, emailOutput.MessageId)
}

func TestSES_ErrCodeMessageRejected(t *testing.T) {
	content := &EmailContent{}
	content.AddRecipientEmail(recipient).AddRecipientEmail("nitin.kothari@continuum.net")
	mockSESClient := &mocks.SESAPI{}
	awsClient := &AWSEmailClient{"us-east-1", nil, nil}
	awsClient.SESClient = mockSESClient
	mockSESClient.On("SendEmail", mock.Anything).Return(nil, awserr.New(ses.ErrCodeMessageRejected, "ErrMessageRejected", errors.New("Error sending emails")))
	_, err := awsClient.SendEmail("", content, awsClient.AWSErrorHandler)
	assert.NotNil(t, err)
}

func TestSES_ErrCodeConfigurationSet(t *testing.T) {
	content := &EmailContent{}
	content.AddRecipientEmail(recipient).AddRecipientEmail("nitin.kothari@continuum.net")
	content.Sender = sender
	mockSESClient := &mocks.SESAPI{}
	awsClient := &AWSEmailClient{"us-east-1", nil, nil}
	awsClient.SESClient = mockSESClient
	mockSESClient.On("SendEmail", mock.Anything).Return(nil, awserr.New(ses.ErrCodeConfigurationSetDoesNotExistException, "ErrConfigurationSetNotSet", errors.New("Error sending emails")))
	_, err := awsClient.SendEmail("", content, awsClient.AWSErrorHandler)
	assert.NotNil(t, err)
}

func TestSES_InvalidReceipients(t *testing.T) {
	content := &EmailContent{}
	content.CharSet = DefaultCharset
	content.TextBody = textBody
	content.Sender = sender
	content.Subject = subject
	content.HTMLBody = htmlBody

	mockSESClient := &mocks.SESAPI{}

	awsClient := &AWSEmailClient{"us-east-1", nil, nil}
	awsClient.SESClient = mockSESClient
	dummyId := "DummyID-1111"
	mockSESClient.On("SendEmail", mock.Anything).Return(&ses.SendEmailOutput{MessageId: &dummyId}, nil)
	_, err := awsClient.SendEmail("", content, awsClient.AWSErrorHandler)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Please provide")
}
