package dto

// ListRoleResponse struct defines response fields
type ListRoleResponse struct {
	Roles []RoleResponse `json:"roles"`
	Total int64          `json:"total"`
}

// RoleResponse struct defines response fields
type RoleResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
	DeletedAt   *string              `json:"deleted_at"`
}
