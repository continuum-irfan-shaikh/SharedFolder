package common

import (
	"testing"
	"time"
)

var testTime = time.Now()

func TestValidTime(t *testing.T) {
	testCases := []struct {
		name     string
		start    time.Time
		end      time.Time
		mustBeOk bool
	}{
		{
			name:     "good case",
			start:    testTime.Add(time.Second),
			end:      testTime,
			mustBeOk: true,
		},
		{
			name:     "bad case",
			start:    testTime,
			end:      testTime.Add(time.Second),
			mustBeOk: false,
		},
	}

	for _, testCase := range testCases {
		test := testCase
		t.Run(test.name, func(t *testing.T) {
			ok := ValidTime(test.start, test.end)
			if ok != test.mustBeOk {
				t.Fatalf("ValidTime: must be %t, but not. Start time = %v, end time = %v", test.mustBeOk, test.start, test.end)
			}
		})
	}
}

func TestCategoryAlreadyExists(t *testing.T) {
	categories := []string{
		"test1",
		"test2",
		"test3",
	}
	testCases := []struct {
		name           string
		category       string
		expectedResult bool
	}{
		{
			name:           "testCase1",
			category:       "test1",
			expectedResult: true,
		},
		{
			name:           "testCase2",
			category:       "test4",
			expectedResult: false,
		},
	}
	for _, testCase := range testCases {
		if result := IsElementAlreadyExists(testCase.category, categories); result != testCase.expectedResult {
			t.Fatalf("%s expect %t, got %t", testCase.name, testCase.expectedResult, result)
		}
	}
}

func TestValidFutureTime(t *testing.T) {
	testCases := []struct {
		name     string
		t        time.Time
		mustBeOk bool
	}{
		{
			name:     "good case",
			t:        time.Now().UTC().Add(5 * time.Minute),
			mustBeOk: true,
		},
		{
			name:     "bad case",
			t:        time.Now().UTC().Add(-5 * time.Minute),
			mustBeOk: false,
		},
	}

	for _, testCase := range testCases {
		test := testCase
		t.Run(test.name, func(t *testing.T) {
			ok := ValidFutureTime(test.t)
			if ok != test.mustBeOk {
				t.Fatalf("ValidFutureTime: must be %t, but not. Time = %v, current UTC time  = %v", test.mustBeOk, test.t, time.Now().UTC())
			}
		})
	}
}
