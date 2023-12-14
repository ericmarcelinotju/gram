package dto

import "time"

// AuditDto struct defines dto for audit entity
type AuditDto struct {
	ID            string    `json:"id"`
	Date          time.Time `json:"date"`
	EntityName    string    `json:"entity_name"`
	EntityId      string    `json:"entity_id"`
	OldValue      string    `json:"old_value"`
	NewValue      string    `json:"new_value"`
	OperationType string    `json:"operation_type"`
	Origin        string    `json:"origin"`

	UserID string  `json:"user_id"`
	User   UserDto `json:"user"`

	PermissionID string        `json:"permission_id"`
	Permission   PermissionDto `json:"permission"`
}
