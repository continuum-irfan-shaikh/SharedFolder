package errorCodePair

import "fmt"

const (
	ERROR = "Err"
	WARN  = "Warn"
)

const (
	//*************FILE SYSTEM ERRORS**********************
	FileSystem = "FileSystem"

	AccessDenied = "AccessDenied"
	Diskfull     = "Diskfull"
	FileNotFound = "FileNotFound"
	Database     = "Database"

	//*************INTERNAL ERRORS*************************
	Internal = "Internal"

	FalseAlarm     = "FalseAlarm"
	ProcessRunning = "ProcessRunning"
	InstallFailure = "InstallFailure"

	//Default Error / Errors due to existing bug
	Operational = "Operational"

	//*************NETWORK ERRORS**************************
	Network = "Network"

	Connection = "Connection"
	Proxy      = "Proxy"

	//*************DOWNLOAD ERROR*************************
	Download = "Download"

	ChecksumValidationFailed = "ChecksumValidationFailed" //File reading and verification errrors after download

)

type errorCode struct {
	Type, App, Error string
}

func newSubErrorCode(err_typ, err_app string) *errorCode {
	return &errorCode{Type: err_typ, App: err_app, Error: Operational}
}

func newMainErrorCode(err_typ, err_app string) *errorCode {
	return &errorCode{Type: err_typ, App: err_app, Error: Internal}
}

func (ec errorCode) String() string {
	return fmt.Sprintf("%v_%v_%v", ec.Type, ec.App, ec.Error)
}

type ErrorCodePair interface {
	SetErrorCodes(main_err_code, sub_err_code string)
	GetErrorDetails() (string, string, string)
}

type errorCodePair struct {
	MainErrorCode, SubErrorCode *errorCode
	ErrorMesssage               string
}

func NewErrorCodePair(err_typ, appName, errorMessage string) ErrorCodePair {
	return &errorCodePair{MainErrorCode: newMainErrorCode(err_typ, appName), SubErrorCode: newSubErrorCode(err_typ, appName), ErrorMesssage: errorMessage}

}

func (ecp *errorCodePair) SetErrorCodes(main_err_code, sub_err_code string) {
	ecp.MainErrorCode.Error, ecp.SubErrorCode.Error = main_err_code, sub_err_code
}

func (ecp errorCodePair) GetErrorDetails() (string, string, string) {
	return ecp.MainErrorCode.String(), ecp.SubErrorCode.String(), ecp.ErrorMesssage
}
