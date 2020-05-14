package zookeeper

import (
	"errors"
	"reflect"
	"testing"

	"github.com/maraino/go-mock"
	"github.com/samuel/go-zookeeper/zk"
)

func TestGetList(t *testing.T) {
	zkMockObj, originalClient := InitMock()
	defer Restore(originalClient)

	expectedResult := []string{"test"}

	zkMockObj.When("Children", mock.Any).Return(expectedResult, &zk.Stat{}, nil)

	arr, err := Queue.GetList("test")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(arr, expectedResult) {
		t.Errorf("expected result: %s, got: %s", expectedResult, arr)
	}
}

func TestCreateItem(t *testing.T) {
	zkMockObj, originalClient := InitMock()
	defer Restore(originalClient)

	expectedResult := "expected result"

	zkMockObj.When("CreateRecursive", mock.Any, mock.Any, mock.Any, mock.Any).Return(expectedResult, nil)

	result, err := Queue.CreateItem([]byte("test"), "test_name")
	if err != nil {
		t.Fatal(err)
	}
	if result != expectedResult {
		t.Errorf("expected result: %s, got: %s", expectedResult, result)
	}
}

func TestGetItemData(t *testing.T) {
	zkMockObj, originalClient := InitMock()
	defer Restore(originalClient)

	expectedResult := []byte("test")

	zkMockObj.When("Get", mock.Any).Return(expectedResult, &zk.Stat{}, nil)

	arr, err := Queue.GetItemData("test", "test")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(arr, expectedResult) {
		t.Errorf("expected result: %s, got: %s", expectedResult, arr)
	}
}

func TestRemoveItem(t *testing.T) {
	zkMockObj, originalClient := InitMock()
	defer Restore(originalClient)

	expectedErr := errors.New("some error")

	zkMockObj.When("Delete", mock.Any, mock.Any).Return(expectedErr)

	err := Queue.RemoveItem("test", "test")
	if err != expectedErr {
		t.Errorf("expected err: %s, got: %s", expectedErr, err)
	}
}
