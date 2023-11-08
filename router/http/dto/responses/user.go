package dto

// ListResponse struct defines response fields
type ListUserResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}

// Response struct defines response fields
type UserResponse struct {
	ID         string       `json:"id"`
	Username   string       `json:"username"`
	Email      string       `json:"email"`
	Firstname  string       `json:"first_name"`
	Lastname   string       `json:"last_name"`
	Department string       `json:"department"`
	Title      string       `json:"title"`
	Avatar     *string      `json:"avatar,omitempty"`
	RoleID     string       `json:"role_id"`
	RoleName   string       `json:"role_name"`
	Role       RoleResponse `json:"role"`

	LastLogin *string `json:"last_login"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
