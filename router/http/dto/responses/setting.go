package dto

// ListSettingResponse struct defines response fields
type ListSettingResponse struct {
	Settings []SettingResponse `json:"settings"`
}

// Response struct defines response fields
type SettingResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
