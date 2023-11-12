package model

import "github.com/ericmarcelinotju/gram/dto"

// SettingEntity struct defines the database model for a setting.
type SettingEntity struct {
	Name  string `gorm:"primaryKey"`
	Value string
}

func NewSettingEntity(entity *dto.SettingDto) *SettingEntity {
	return &SettingEntity{
		Name:  entity.Name,
		Value: entity.Value,
	}
}

func (entity *SettingEntity) ToDto() *dto.SettingDto {
	return &dto.SettingDto{
		Name:  entity.Name,
		Value: entity.Value,
	}
}
