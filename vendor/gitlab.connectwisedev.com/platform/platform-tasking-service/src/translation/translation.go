package translation

import (
	"path"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation/genlang"
)

const (
	fallbackLanguage = "en-US"
	assetFolder      = "translation/languages"
	// ErrTranslator Error message returned in case translator can't be created
	ErrTranslator       = "can't create translator"
	ErrorCodeTranslator = `error_cannot_create_translator`
)

var (
	defaultBundle        = bundle.New()
	defaultTranslateFunc bundle.TranslateFunc
	mockedTranslation    = map[string]interface{}{
		"id":          "place_holder",
		"translation": "foobar",
	}
)

// TranslatorType translator's structure
type TranslatorType struct {
	Language      string
	translateFunc bundle.TranslateFunc
}

// Load loads config with chosen language
func Load() error {
	err := loadFiles(assetFolder)
	if err != nil {
		return err
	}

	if config.Config.DefaultLanguage == "" {
		config.Config.DefaultLanguage = fallbackLanguage
	}

	defaultTranslateFunc, err = defaultBundle.Tfunc(config.Config.DefaultLanguage)
	return err
}

// New creates new translator
func New(language string) (result TranslatorType, err error) {
	translationFunc, parsedLanguage, err := defaultBundle.TfuncAndLanguage(language, config.Config.DefaultLanguage)
	if err != nil {
		return
	}
	result = TranslatorType{
		Language:      parsedLanguage.String(),
		translateFunc: translationFunc,
	}
	return
}

// MockTranslations for tests
func MockTranslations() error {
	config.Config.DefaultLanguage = "en-US"

	t, err := translation.NewTranslation(mockedTranslation)
	if err != nil {
		return err
	}

	defaultBundle.AddTranslation(language.Parse(config.Config.DefaultLanguage)[0], t)

	defaultTranslateFunc = func(translationID string, _ ...interface{}) string {
		return translationID
	}

	return nil
}

// Translate returns t
//he message in chosen language
func (t TranslatorType) Translate(id string, args ...interface{}) (result string) {
	result = t.translateFunc(id, args...)
	if result == id {
		result = defaultTranslateFunc(id, args...)
	}
	return
}

func loadFiles(folder string) error {
	files, err := genlang.AssetDir(folder)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := path.Join(folder, file)
		data, err := genlang.Asset(fullPath)
		if err != nil {
			return err
		}

		err = defaultBundle.ParseTranslationFileBytes(fullPath, data)
		if err != nil {
			return err
		}
	}

	return nil
}
