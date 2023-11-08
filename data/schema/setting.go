package schema

import (
	"github.com/ericmarcelinotju/gram/domain/model"
)

// Setting struct defines the database model for a setting.
type Setting struct {
	Name  string `gorm:"primaryKey"`
	Value string
}

func NewSettingSchema(entity *model.Setting) *Setting {
	return &Setting{
		Name:  entity.Name,
		Value: entity.Value,
	}
}

func (entity *Setting) ToDomainModel() *model.Setting {
	return &model.Setting{
		Name:  entity.Name,
		Value: entity.Value,
	}
}
