package translation

import (
	"fmt"
	"os"
	"testing"

	rLog "gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
)

const (
	error400En = "Can not decode input data"
	error400Uk = "Неможливо декодувати вхідні дані"
)

var (
	emptyIf = []interface{}{}
)

func init() {
	loggerForTest()
}

func Test_Load(t *testing.T) {
	c := ErrorCodeTranslator
	c = fallbackLanguage
	c = assetFolder
	c = ErrTranslator
	fmt.Println(c)

	err := Load()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Translate(t *testing.T) {
	config.Config.DefaultLanguage = "en-US"
	Load()

	var testCases = []struct {
		language string
		id       string
		args     interface{}
		expected string
	}{
		{"en-US", "error_cant_decode_input_data", emptyIf, error400En},
		{"en", "error_cant_decode_input_data", emptyIf, error400En},
		{"fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5", "error_cant_decode_input_data", emptyIf, error400En},
		{"", "error_cant_decode_input_data", emptyIf, error400En},
		{"zz", "error_cant_decode_input_data", emptyIf, error400En},
		{"uk", "error_cant_decode_input_data", emptyIf, error400Uk},
		{"en-US", "zzz_999999", emptyIf, "zzz_999999"},
		{"fr-CH, fr;q=0.9, uk;q=0.8, de;q=0.7, *;q=0.5", "error_cant_decode_input_data", emptyIf, error400Uk},
		{"uk-UA", "error_cant_decode_input_data", map[string]interface{}{"ID": "UUID1"}, error400Uk},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test_case_%d", i), func(t *testing.T) {
			translator, err := New(testCase.language)
			if err != nil {
				t.Error(err)
			}

			got := translator.Translate(testCase.id, testCase.args)
			if got != testCase.expected {
				t.Errorf("expected: '%s', but got: '%s'", testCase.expected, got)
			}
		})
	}
}

func loggerForTest() {
	config.Config.Log.FileName = "logs_test.log"
	config.Config.Log.LogLevel = rLog.INFO
	config.Config.DefaultLanguage = "en-US"
	logger.Load(config.Config.Log)
	Load()
	defer os.Remove(config.Config.Log.FileName)
}

func TestLoad(t *testing.T) {
	config.Config.DefaultLanguage = ""
	err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if config.Config.DefaultLanguage != fallbackLanguage {
		t.Fatalf("Wrong language loaded\n")
	}
}

func Test_MockTranslations(t *testing.T) {
	expected := mockedTranslation["translation"].(string)

	t.Run("Test_mock", func(t *testing.T) {
		err := MockTranslations()
		if err != nil {
			t.Error(err)
		}
	})

	if defaultTranslateFunc("hello") != "hello" {
		t.Fatalf("Something went wrong")
	}

	t.Run("Test_translation", func(t *testing.T) {
		translator, _ := New(fallbackLanguage)

		got := translator.Translate(mockedTranslation["id"].(string), emptyIf)
		if got != expected {
			t.Errorf("expected: '%s', but got: '%s'", expected, got)
		}
	})
}

func TestNew(t *testing.T) {
	lang := "en-us"
	translatorTest, err := New(lang)
	if translatorTest.Language != lang {
		t.Fatalf("Got %s, want %s", translatorTest.Language, lang)
	}

	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadFiles(t *testing.T) {
	folder := "/home"
	err := loadFiles(folder)
	if err == nil {
		t.Fatalf("Error expected, but no error returned")
	}
}
