package currency

import (
	"errors"
	"strings"

	"github.com/leekchan/accounting"
)

type Currency struct {
	Code       string
	Locale     accounting.Locale
	Accounting accounting.Accounting
}

func NewCurrency(code string) (*Currency, error) {
	locale, ok := accounting.LocaleInfo[strings.ToUpper(code)]
	if !ok {
		return nil, errors.New("invalid currency code")
	}
	return &Currency{
		Code:   code,
		Locale: locale,
		Accounting: accounting.Accounting{
			Symbol:    locale.ComSymbol,
			Precision: locale.FractionLength,
			Thousand:  locale.ThouSep,
			Decimal:   locale.DecSep,
		},
	}, nil
}

func (c Currency) Format(value float64) string {
	return c.Accounting.FormatMoney(value)
}
