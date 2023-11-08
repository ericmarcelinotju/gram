package dto

// ListResponse struct defines response fields
type ListAuditResponse struct {
	Audits []AuditResponse `json:"audits"`
	Total  int64           `json:"total"`
}

// Response struct defines response fields
type AuditResponse struct {
	ID            string `json:"id"`
	Date          string `json:"date"`
	EntityName    string `json:"entity_name"`
	EntityId      string `json:"entity_id"`
	OldValue      string `json:"old_value"`
	NewValue      string `json:"new_value"`
	OperationType string `json:"operation_type"`
	Origin        string `json:"origin"`

	UserID string       `json:"user_id"`
	User   UserResponse `json:"user"`

	PermissionID string             `json:"permission_id"`
	Permission   PermissionResponse `json:"permission"`
}
