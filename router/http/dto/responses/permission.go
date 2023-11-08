package dto

// ListResponse struct defines response fields
type ListPermissionResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
	Total       int64                `json:"total"`
}

// Response struct defines response fields
type PermissionResponse struct {
	ID          string  `json:"id"`
	Method      string  `json:"method"`
	Module      string  `json:"module"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at"`
}
