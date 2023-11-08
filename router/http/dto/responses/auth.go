package dto

// LoginResponse struct defines response fields
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
