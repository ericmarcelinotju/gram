package dto

// SettingDto struct defines dto of setting entity
type SettingDto struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PostSettingDto struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}
