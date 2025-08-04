package localization

import (
	"github.com/ericmarcelinotju/gram/localization"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Translate(str string, lang string, stringMap map[string]string) (string, error) {
	localizer := i18n.NewLocalizer(localization.Get().Bundle, lang) //3
	localizeConfig := i18n.LocalizeConfig{                          //2
		MessageID:    str,
		TemplateData: stringMap,
	}

	return localizer.Localize(&localizeConfig) //3

}
