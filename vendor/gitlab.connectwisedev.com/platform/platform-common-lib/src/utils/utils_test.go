package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	errorCodes "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/errorCodePair"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/constants"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/plugin/protocol"
)

func TestToString(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected string
	}{
		{nil, ""},
		{"a", "a"},
		{1, ""},
	}
	for _, d := range testCases {
		returned := ToString(d.v)
		if returned != d.expected {
			t.Errorf("Unexpected object/value in converting interface to string")
		}
	}
}

func TestInt64SliceToStringSlice(t *testing.T) {
	testCases := []struct {
		v        []int64
		expected []string
	}{
		{[]int64{}, nil},
		{[]int64{1}, []string{"1"}},
		{[]int64{1, 2, 4, 5}, []string{"1", "2", "4", "5"}},
		{[]int64{0, 1, 2, 4, 5}, []string{"0", "1", "2", "4", "5"}},
	}
	
	for _, d := range testCases {
		returned := Int64SliceToStringSlice(d.v)
		if !reflect.DeepEqual(returned, d.expected) {
			t.Errorf("Unexpected object/value in converting interface to string array")
			break
		}
	}
}

func TestToInt64(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected int64
	}{
		{nil, 0},
		{"a", 0},
		{int64(1), 1}, //int can not be type casted to int64 that is why it needs to be converted to int64
	}
	for _, d := range testCases {
		returned := ToInt64(d.v)
		if returned != d.expected {
			t.Errorf("Unexpected object/value in converting interface to int64")
		}
	}
}

func TestToInt(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected int
	}{
		{nil, 0},
		{"a", 0},
		{1, 1},
	}
	for _, d := range testCases {
		returned := ToInt(d.v)
		if returned != d.expected {
			t.Errorf("Unexpected object/value in converting interface to int")
		}
	}
}

func TestToStringArray(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected []string
	}{
		{nil, []string{}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{0, []string{}},
	}
	for _, d := range testCases {
		returned := ToStringArray(d.v)
		if !reflect.DeepEqual(returned, d.expected) {
			t.Errorf("Unexpected object/value in converting interface to string array")
			break
		}
	}
}

func TestToByteArray(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected []byte
	}{
		{nil, []byte{}},
		{[]byte("ABC"), []byte("ABC")},
		{0, []byte{}},
	}
	for _, d := range testCases {
		returned := ToByteArray(d.v)
		if !reflect.DeepEqual(returned, d.expected) {
			t.Errorf("Unexpected object/value in converting interface to string array")
			break
		}
	}
}

func TestGetTransactionID(t *testing.T) {
	tid := GetTransactionID()
	//check if tid is a valid uuid
	_, err := uuid.Parse(tid)
	if err != nil {
		t.Errorf("Not a valid uuid returned %s; error %v", tid, err)
	}
}

func TestGetTransactionIDFromResponse(t *testing.T) {
	id := "1"
	r := &http.Response{}
	r.Header = make(http.Header, 1)
	r.Header.Set(string(protocol.HdrTransactionID), id)

	tid := GetTransactionIDFromResponse(r)
	if tid != id {
		t.Errorf("Unexpected transactionId returned , Expected:%s, Returned:%s", id, tid)
	}

}

func TestGetTransactionIDFromResponseNilCheck(t *testing.T) {

	tid := GetTransactionIDFromResponse(nil)
	if tid != "" {
		t.Errorf("Unexpected transactionId returned , Returned:%s", tid)
	}

}

func TestGetTransactionIDFromRequest(t *testing.T) {
	t.Run("1. Transaction ID exists in protocol.HdrTransactionID", func(t *testing.T) {
		id := "1"
		r := &http.Request{}
		r.Header = make(http.Header, 1)
		r.Header.Set(string(protocol.HdrTransactionID), id)

		tid := GetTransactionIDFromRequest(r)
		if tid != id {
			t.Errorf("Unexpected transactionId returned , Expected:%s, Returned:%s", id, tid)
		}
	})
	t.Run("2. Transaction ID exists in constants.TransactionID", func(t *testing.T) {
		id := "1"
		r := &http.Request{}
		r.Header = make(http.Header, 1)
		r.Header.Set(constants.TransactionID, id)

		tid := GetTransactionIDFromRequest(r)
		if tid != id {
			t.Errorf("Unexpected transactionId returned , Expected:%s, Returned:%s", id, tid)
		}
	})

	t.Run("3. Transaction ID does not exists in any header, return uuid", func(t *testing.T) {
		r := &http.Request{}
		tid := GetTransactionIDFromRequest(r)
		//check if tid is a valid uuid
		_, err := uuid.Parse(tid)
		if err != nil {
			t.Errorf("Not a valid uuid returned %s in case of no transaction header present; error %v", tid, err)
		}
	})

	t.Run("4. Transaction ID exists in both header, return X-RequestID", func(t *testing.T) {
		requestID := "1"
		transactionID := "2"
		r := &http.Request{}
		r.Header = make(http.Header, 2)
		r.Header.Set(string(protocol.HdrTransactionID), requestID)
		r.Header.Set(constants.TransactionID, transactionID)

		tid := GetTransactionIDFromRequest(r)
		if tid != requestID {
			t.Errorf("Unexpected transactionId returned , Expected:%s, Returned:%s", requestID, tid)
		}
	})

	t.Run("5. Request is nil, return uuid", func(t *testing.T) {
		tid := GetTransactionIDFromRequest(nil)
		//check if tid is a valid uuid
		_, err := uuid.Parse(tid)
		if err != nil {
			t.Errorf("Not a valid uuid returned %s in case of no transaction header present; error %v", tid, err)
		}
	})
}

func TestGetChecksumFromRequest(t *testing.T) {
	tid := GetChecksumFromRequest(nil)
	if tid != "" {
		t.Errorf("Unexpected transactionId returned , Returned:%s", tid)
	}

}

func TestGetChecksum(t *testing.T) {
	validChecksum := "0cbc6611f5540bd0809a388dc95a615b"
	message := []byte("Test")
	checksum := GetChecksum(message)
	if checksum != validChecksum {
		t.Errorf("Failed!:%v", checksum)
	}
}

func TestGetChecksumBlank(t *testing.T) {
	validChecksum := "d41d8cd98f00b204e9800998ecf8427e"
	message := []byte{}
	checksum := GetChecksum(message)
	if checksum != validChecksum {
		t.Errorf("Failed!: Unexpected checksum %s", checksum)
	}
}

func TestValidateMessageFail(t *testing.T) {
	checksum := "TestString"
	flag, _ := ValidateMessage([]byte("TestString"), checksum)
	if flag {
		t.Errorf("Failed!: Function should return false when checksum is invalid")
	}
}

func TestValidateMessageSuccess(t *testing.T) {
	message := []byte("TestString")
	checksum := GetChecksum(message)
	flag, _ := ValidateMessage(message, checksum)
	if !flag {
		t.Errorf("Failed!: Function should return true when checksum %s", checksum)
	}
}

func TestValidateMessageBlankCheckSum(t *testing.T) {
	message := []byte("TestString")
	flag, _ := ValidateMessage(message, "")
	if !flag {
		t.Errorf("Failed!: Function should return true when checksum blank")
	}
}

func systemCall() (interface{}, error) {
	time.Sleep(3 * time.Second)
	return "Success", nil
}

func systemCallTimeoutHandler() {
	fmt.Println("In systemCallTimeoutHandler method")
}

type RecoverTest struct {
	Panic bool
	T     *testing.T
}

func (r RecoverTest) PanicHandler(err error) {
	if r.Panic && err == nil {
		r.T.Errorf("Expected Error but Got Nil")
	}

	if !r.Panic && err != nil {
		r.T.Errorf("Expected nil but Got Error %v", err)
	}
}

func (r RecoverTest) PanicGenerator() {
	if r.Panic {
		panic(fmt.Errorf("Generating Panic"))
	}
}

func TestToBool(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TC1_True",
			args: args{
				v: true,
			},
			want: true,
		},
		{
			name: "TC2_NIL",
			args: args{
				v: nil,
			},
			want: false,
		},
		{
			name: "TC3_Int",
			args: args{
				v: 1,
			},
			want: false,
		},
		{
			name: "TC4_False",
			args: args{
				v: false,
			},
			want: false,
		},
		{
			name: "TC5_Float64",
			args: args{
				v: float64(1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBool(tt.args.v); got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "TC1_String",
			args: args{
				v: "abc",
			},
			want: 0,
		},
		{
			name: "TC2_Nil",
			args: args{
				v: nil,
			},
			want: 0,
		},
		{
			name: "TC3_float64",
			args: args{
				v: float64(1),
			},
			want: 1,
		},
		{
			name: "TC4_false",
			args: args{
				v: false,
			},
			want: 0,
		},
		{
			name: "TC5_true",
			args: args{
				v: true,
			},
			want: 0,
		},
		{
			name: "TC6_true",
			args: args{
				v: int64(23),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFloat64(tt.args.v); got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStringMap(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "TC1",
			args: args{v: nil},
			want: nil,
		},
		{
			name: "TC2",
			args: args{v: map[string]string{
				"name":  "a",
				"value": "a",
			},
			},
			want: map[string]string{
				"name":  "a",
				"value": "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStringMap(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetermineErrorCodePair(t *testing.T) {
	type args struct {
		errMsg string
	}
	tests := []struct {
		name          string
		args          args
		wantMainError string
		wantSubError  string
	}{
		{
			name:          "When error category is File System",
			args:          args{"FileSystem|FileNotFound"},
			wantMainError: errorCodes.FileSystem,
			wantSubError:  errorCodes.FileNotFound,
		},
		{
			name:          "When error category is Network",
			args:          args{"Network|Proxy"},
			wantMainError: errorCodes.Network,
			wantSubError:  errorCodes.Proxy,
		},
		{
			name:          "When error category is Download",
			args:          args{"Download|ChecksumValidationFailed"},
			wantMainError: errorCodes.Download,
			wantSubError:  errorCodes.ChecksumValidationFailed,
		},
		{
			name:          "When error category is Internal",
			args:          args{"Internal|InstallFailure"},
			wantMainError: errorCodes.Internal,
			wantSubError:  errorCodes.InstallFailure,
		},
		{
			name:          "When error category is uncategorized",
			args:          args{"uncategorized"},
			wantMainError: errorCodes.Internal,
			wantSubError:  errorCodes.Operational,
		},
		{
			name:          "When error category is malformed",
			args:          args{"Internal|AccessDenied"},
			wantMainError: errorCodes.Internal,
			wantSubError:  errorCodes.Operational,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMainError, gotSubError := DetermineErrorCodePair(tt.args.errMsg)
			if gotMainError != tt.wantMainError {
				t.Errorf("DetermineErrorCodePair() gotMainError = %v, want %v", gotMainError, tt.wantMainError)
			}
			if gotSubError != tt.wantSubError {
				t.Errorf("DetermineErrorCodePair() gotSubError = %v, want %v", gotSubError, tt.wantSubError)
			}
		})
	}
}

func Test_EqualFold(t *testing.T) {
	tests := []struct {
		name   string
		source []string
		target []string
		want   bool
	}{
		{
			name:   "single string, matching, with same case",
			source: []string{"abc"}, target: []string{"abc"}, want: true,
		},
		{
			name:   "single string, matching, with case difference",
			source: []string{"abc"}, target: []string{"Abc"},
			want: true,
		},
		{
			name:   "single string alphanumeric, matching, with case difference",
			source: []string{"ab1c"}, target: []string{"Ab1c"},
			want: true,
		},
		{
			name:   "multiple strings,already sorted, matching, with same case",
			source: []string{"abc", "xyz"}, target: []string{"abc", "xyz"},
			want: true,
		},
		{
			name:   "multiple strings,without sorted, matching, with same case",
			source: []string{"abc", "xyz"}, target: []string{"xyz", "abc"},
			want: true,
		},
		{
			name:   "multiple strings,already sorted, matching, with case difference",
			source: []string{"abc", "xyz"}, target: []string{"aBc", "XYZ"},
			want: true,
		},
		{
			name:   "multiple strings,without sorted, matching, with case difference",
			source: []string{"aBC12", "xyz"}, target: []string{"xyZ", "abc12"},
			want: true,
		},
		{
			name:   "multiple strings,not matching, case 1",
			source: []string{"abcmno", "xyz"}, target: []string{"xyz", "abc"},
			want: false,
		},
		{
			name:   "multiple strings,length mismatch",
			source: []string{"abc", "xyz"}, target: []string{"xyz"},
			want: false,
		},
		{
			name:   "single string, not matching, case 2",
			source: []string{"abc"}, target: []string{"Abc123"},
			want: false,
		},
		{
			name:   "single string, not matching, case 3",
			source: []string{"abc"}, target: []string{"1abc"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualFold(tt.source, tt.target); got != tt.want {
				t.Errorf("EqualFold() = %v, want %v", got, tt.want)
			}
		})
	}
}
