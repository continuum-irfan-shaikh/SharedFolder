package utils

//Package utils is to convert interface type to specific type

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"crypto/md5"
	"encoding/hex"

	errorCodes "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/errorCodePair"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/constants"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/plugin/protocol"
	"github.com/google/uuid"
)

//ToString converts an interface type to string
func ToString(v interface{}) string {
	t, ok := v.(string)
	if ok {
		return t
	}
	return ""
}

//ToTime converts an interface type to time
func ToTime(v interface{}) time.Time {
	t, ok := v.(time.Time)
	if ok {
		return t
	}
	return time.Time{}
}

//ToInt64 converts an interface type to int64
//interface{} holding an int will not be type casted to int64 and will return 0 as the result
func ToInt64(v interface{}) int64 {
	t, ok := v.(int64)
	if ok {
		return t
	}
	return 0
}

//ToInt converts an interface type to int
func ToInt(v interface{}) int {
	t, ok := v.(int)
	if ok {
		return t
	}
	return 0
}

//ToFloat64 converts an interface type to float64
func ToFloat64(v interface{}) float64 {
	t, ok := v.(float64)
	if ok {
		return t
	}
	return 0
}

//ToBool converts an interface type to bool
func ToBool(v interface{}) bool {
	t, ok := v.(bool)
	if ok {
		return t
	}
	return false
}

//ToStringArray converts an interface type to string array
func ToStringArray(v interface{}) []string {
	t, ok := v.([]string)
	if ok {
		return t
	}
	return []string{}
}

//ToByteArray converts an interface type to byte array
func ToByteArray(v interface{}) []byte {
	t, ok := v.([]byte)
	if ok {
		return t
	}
	return []byte{}
}

//ToStringMap converts an interface type to map[string]string
func ToStringMap(v interface{}) map[string]string {
	t, ok := v.(map[string]string)
	if ok {
		return t
	}
	return nil
}

//GetTransactionID generates new transactionid
func GetTransactionID() string {
	return uuid.New().String()
}

//GetTransactionIDFromResponse retrieves transactionid from the Response header
func GetTransactionIDFromResponse(res *http.Response) string {
	if res == nil {
		return ""
	}
	return res.Header.Get(string(protocol.HdrTransactionID))
}

//GetTransactionIDFromRequest retrieves transactionID from the Request header
func GetTransactionIDFromRequest(req *http.Request) string {
	if req == nil {
		return GetTransactionID()
	}
	value := GetValueFromRequestHeader(req, protocol.HdrTransactionID)
	if value != "" {
		return value

	}
	value = req.Header.Get(constants.TransactionID)
	if value != "" {
		return value
	}
	return GetTransactionID()
}

//GetQueryValuesFromRequest to get query values from request for given filter
func GetQueryValuesFromRequest(req *http.Request, filter string) []string {
	queryValues := req.URL.Query()
	if _, ok := queryValues[filter]; ok {
		return queryValues[filter]
	}
	return []string{}
}

//GetChecksumFromRequest retrives MD5 from the request header
func GetChecksumFromRequest(req *http.Request) string {
	return GetValueFromRequestHeader(req, protocol.HdrContentMD5)
}

//GetValueFromRequestHeader retrieves header value for given Key from the Request header
func GetValueFromRequestHeader(req *http.Request, key protocol.HeaderKey) string {
	if req == nil {
		return ""
	}

	return req.Header.Get(string(key))
}

//GetChecksum is a function to calculate MD5 hash value
func GetChecksum(message []byte) string {
	hasher := md5.New()
	hasher.Write(message) //nolint
	return hex.EncodeToString(hasher.Sum(nil))
}

//ValidateMessage checks if message is corrupted or not
func ValidateMessage(message []byte, receievedChecksum string) (bool, string) {
	hashValue := GetChecksum(message)
	if receievedChecksum != "" && hashValue != receievedChecksum {
		return false, hashValue
	}
	return true, hashValue
}

// changes for agent autoupdate error standardization START HERE. To be refactored as per common-lib standards for comming rollouts

// Determine Error code pairs for autoupdate failures
func DetermineErrorCodePair(errMsg string) (mainError, subError string) {
	//determine errors code pairs
	if strings.Contains(errMsg, errorCodes.FileSystem) {
		mainError, subError = errorCodes.FileSystem, determineFileSystemErrorCodes(errMsg)
	} else if strings.Contains(errMsg, errorCodes.Network) {
		mainError, subError = errorCodes.Network, determineNetworkErrorCodes(errMsg)
	} else if strings.Contains(errMsg, errorCodes.Download) {
		mainError, subError = errorCodes.Download, determineDownloadErrorCodes(errMsg)
	} else if strings.Contains(errMsg, errorCodes.Internal) {
		mainError, subError = errorCodes.Internal, determineInternalErrorCodes(errMsg)
	}
	if subError == "" {
		mainError, subError = errorCodes.Internal, errorCodes.Operational
	}
	return
}

func determineFileSystemErrorCodes(errMsg string) string {
	if strings.Contains(errMsg, errorCodes.Database) {
		return errorCodes.Database
	} else if strings.Contains(errMsg, errorCodes.AccessDenied) {
		return errorCodes.AccessDenied
	} else if strings.Contains(errMsg, errorCodes.Diskfull) {
		return errorCodes.Diskfull
	} else if strings.Contains(errMsg, errorCodes.FileNotFound) {
		return errorCodes.FileNotFound
	}
	return ""
}

func determineNetworkErrorCodes(errMsg string) string {
	if strings.Contains(errMsg, errorCodes.Connection) {
		return errorCodes.Connection
	} else if strings.Contains(errMsg, errorCodes.Proxy) {
		return errorCodes.Proxy
	}
	return ""
}

func determineDownloadErrorCodes(errMsg string) string {
	if strings.Contains(errMsg, errorCodes.ChecksumValidationFailed) {
		return errorCodes.ChecksumValidationFailed
	}
	return ""
}

func determineInternalErrorCodes(errMsg string) string {
	if strings.Contains(errMsg, errorCodes.FalseAlarm) {
		return errorCodes.FalseAlarm
	} else if strings.Contains(errMsg, errorCodes.ProcessRunning) {
		return errorCodes.ProcessRunning
	} else if strings.Contains(errMsg, errorCodes.InstallFailure) {
		return errorCodes.InstallFailure
	}
	return ""
}

// changes for agent autoupdate error standardization END HERE.

//EqualFold is to check if two strings array are same with case insensitive
func EqualFold(source, target []string) bool {
	sort.Slice(source, func(i, j int) bool { return strings.ToLower(source[i]) < strings.ToLower(source[j]) })
	sort.Slice(target, func(i, j int) bool { return strings.ToLower(target[i]) < strings.ToLower(target[j]) })

	lSource := len(source)
	lTarget := len(target)

	if lSource != lTarget {
		return false
	}

	for i := 0; i < lSource; i++ {
		if !strings.EqualFold(source[i], target[i]) {
			return false
		}
	}

	return true
}

// Int64SliceToStringSlice converts int64 slice to string slice
func Int64SliceToStringSlice(in []int64) (out []string) {
	for _, v := range in {
		s := strconv.Itoa(int(v))
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
