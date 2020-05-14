package gateway

import "time"

//PackageDetails gets the mailbox messages
type PackageDetails struct {
	Checksum      string `json:"checksum,omitempty" gorm:"primary_key"`
	HashAlgorithm string `json:"hashAlgorithm,omitempty"`

	RequestedURL    string `json:"requestedURL,omitempty"`
	Filename        string `json:"filename,omitempty"`
	FileSizeInBytes int64  `json:"fileSizeInBytes,omitempty"`
	FileLocation    string `json:"fileLocation,omitempty"`

	DownloadTimeStampUTC     time.Time `json:"downloadTimeStampUTC,omitempty"`
	LastAccessedTimeStampUTC time.Time `json:"lastAccessedTimeStampUTC,omitempty"`

	Status             string `json:"status,omitempty"`
	DownloadRetryCount int    `json:"downloadRetryCount,omitempty"`
	ServeCount         int    `json:"serveCount,omitempty"`
	Deleted            bool   `json:"deleted,omitempty"`

	ErrorCode    string `json:"errorCode,omitempty"`
	HTTPStatus   int    `json:"httpStatus,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`

	Headers string `json:"headers,omitempty"`
}

const (
	//Blacklisted   Url Blacklisted , due to 5xx erros, checksum mismatch
	Blacklisted = "Blacklisted"
	//Downloading url is currently being downloaded
	Downloading = "Downloading"
	//Downloaded url is downloaded
	Downloaded = "Downloaded"
	// DoesNotExists Package doees not exists in the record
	DoesNotExists = "DoesNotExists"
	// ChecksumValidationFailed failred to validate checksum
	ChecksumValidationFailed = "ChecksumValidationFailed"
)

// PackageStatus ...
type PackageStatus struct {
	Status  string `json:"status,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

//Transaction ...
type Transaction struct {
	TransactionID string `json:"transactionID"`
}
