package config

import (
	"fmt"
	"testing"
)

func TestIsMissingTrue(t *testing.T) {
	conf := Configuration{}
	testVal := conf.isMissing()
	if !testVal {
		t.Error("isMissing don't validates")
	}
}

func TestIsMissingFalse(t *testing.T) {
	conf := Configuration{}
	conf.Log.FileName = "TestLog.log"
	conf.ListenURL = ":0000"
	conf.CassandraURL = "SomeTestURL"
	testVal := conf.isMissing()
	if testVal {
		t.Error("isMissing don't validates")
	}
}

func TestReadConfigFromJSON(t *testing.T) {
	configFilePath := "../../config.json"
	err := readConfigFromJSON(configFilePath)
	if err != nil {
		t.Errorf("Reading configuration failed: %s\n", err)
	}

	configFilePath = "test.json"
	err = readConfigFromJSON(configFilePath)
	if err == nil {
		t.Errorf("Reading configuration failed: %s\n", err)
	}
}

func TestReadConfigFromJSONWrongPath(t *testing.T) {
	configFilePath := ""
	err := readConfigFromJSON(configFilePath)
	if err == nil {
		t.Errorf("Something's going wrong: %s\n", err)
	}
}

func TestReadConfigFromENV(t *testing.T) {
	err := readConfigFromENV()
	if err != nil {
		t.Error("Configuration is missing")
	}
}

func TestLoadReadingFromENV(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error reading from ENV: %s", r)
		}

	}()

	ConfigFilePath := ""
	fmt.Printf("Reading from ENV with epty path (%s)", ConfigFilePath)

	Load()
}

func TestLoadReadingFromJSON(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error reading from JSON: %s", r)
		}

	}()

	ConfigFilePath := "../config.json"
	fmt.Printf("Reading from JSON (%s)", ConfigFilePath)

	Load()
}

//For test coverage
func TestScheduledJob_Geters(t *testing.T) {
	expectedString := "string"
	scheduledJob := ScheduledJob{Name: expectedString, Schedule: expectedString, Task: expectedString}

	gotName := scheduledJob.GetName()
	if gotName != expectedString {
		t.Errorf("Got %v but want %v", gotName, expectedString)
	}
	gotSchedule := scheduledJob.GetSchedule()
	if gotName != expectedString {
		t.Errorf("Got %v but want %v", gotSchedule, expectedString)
	}
	gotTask := scheduledJob.GetTask()
	if gotName != expectedString {
		t.Errorf("Got %v but want %v", gotTask, expectedString)
	}
}
