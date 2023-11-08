package model

import dto "gitlab.com/firelogik/helios/router/http/dto/responses"

// Setting struct defines the database model for a setting.
type Setting struct {
	ID    string
	Name  string
	Value string
}

func (entity *Setting) ToResponseModel() *dto.SettingResponse {
	return &dto.SettingResponse{
		Name:  entity.Name,
		Value: entity.Value,
	}
}
