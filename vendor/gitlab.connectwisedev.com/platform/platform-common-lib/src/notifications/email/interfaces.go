package email

//Enum representation of placeholders
type PlaceholderKey string

//Wrapper interface for any EmailClients
type EmailClient interface {
	SendEmail(transactionID string, emailContent *EmailContent, errorHandler ErrorHandler) (*EmailOutput, error)
}

//Callback error handler invoked in SendEmail
type ErrorHandler func(transactionID string, err error)

//Represents Output type after send email is invoked
type EmailOutput struct {
	MessageId *string
}

//Email Container to be composed. Email contents are built using this container.
type EmailContent struct {
	Sender          string
	ToAddresses     []*string
	CCAddresses     []*string
	Subject         string
	HTMLBody        string
	HTMLTemplate    string
	TextBody        string
	CharSet         string
	ContentKeyValue map[PlaceholderKey]string
	BodyKeyValue    map[PlaceholderKey]string
}
