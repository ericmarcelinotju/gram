package localization

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"path/filepath"
	"sync"
)

type Localization struct {
	Bundle *i18n.Bundle
}

func NewI18n() *Localization {

	//i18n
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	registeredLanguage := []string{
		language.Japanese.String(),
		language.English.String(),
	}
	for _, lang := range registeredLanguage {
		path, err := filepath.Abs("lang/" + lang + ".json")
		if err != nil {
			log.Println(err.Error())
		}
		_, err = bundle.LoadMessageFile(path)
		if err != nil {
			log.Println(err.Error())
		}
	}
	i18n.NewLocalizer(bundle, language.English.String(), language.Japanese.String())
	instance := &Localization{
		Bundle: bundle,
	}
	return instance

}

var instance *Localization
var once sync.Once

func Get() *Localization {
	once.Do(func() {
		instance = NewI18n()
	})

	return instance
}
