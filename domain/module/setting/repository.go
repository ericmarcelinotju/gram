package setting

import (
	"context"

	"github.com/ericmarcelinotju/gram/domain/model"
)

// Repository provides an abstraction on top of the setting data source
type Repository interface {
	SaveSetting(context.Context, string, string) error
	SelectSetting(context.Context) ([]model.Setting, error)
	SelectSettingByName(context.Context, string) (string, error)
	DeleteSetting(context.Context, string) error
}
