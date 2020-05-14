package mocks

import (
	"encoding/json"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
)

const cantDecodeInputData = "Can not decode input data"

func TestIsErrorMessageValid(t *testing.T) {
	translation.Load()

	type Message struct {
		Message string
	}

	cannotDecodeInputDataMessage, err := json.Marshal(Message{Message: cantDecodeInputData})
	if err != nil {
		t.Errorf("Cannot decode errorcode constant")
	}

	wrongMessage, err := json.Marshal(Message{Message: "The other message"})
	if err != nil {
		t.Errorf("Cannot decode errorcode constant")
	}

	testCases := []struct {
		name                string
		inputExpectedErrMsg string
		inputReceivedErrMsg string
		expectedBool        bool
	}{
		{
			name:                "testCase 0 - cannot parse received json",
			expectedBool:        false,
			inputExpectedErrMsg: errorcode.CodeUpdated,
			inputReceivedErrMsg: "{invalidJsonToParse",
		},
		{
			name:                "testCase 1 - cannot parse received message",
			expectedBool:        false,
			inputExpectedErrMsg: errorcode.ErrorCantDecodeInputData,
			inputReceivedErrMsg: string(wrongMessage),
		},
		{
			name:                "testCase 2 - Ok",
			expectedBool:        true,
			inputExpectedErrMsg: errorcode.ErrorCantDecodeInputData,
			inputReceivedErrMsg: string(cannotDecodeInputDataMessage),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := IsErrorMessageValid(tc.inputExpectedErrMsg, []byte(tc.inputReceivedErrMsg))
			if got != tc.expectedBool {
				t.Errorf("Wanted %v but got %v", tc.expectedBool, got)
			}
		})
	}
}
