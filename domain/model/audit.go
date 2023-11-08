package model

import (
	"time"

	dto "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	timeUtils "github.com/ericmarcelinotju/gram/utils/time"
)

// Audit struct defines the database model for a Audit.
type Audit struct {
	ID            string
	Date          time.Time
	EntityName    string
	EntityId      string
	OldValue      string
	NewValue      string
	OperationType string
	Origin        string

	UserID string
	User   User

	PermissionID string
	Permission   Permission

	Pagination Pagination
	Sort       Sort
}

func (entity *Audit) ToResponseModel() *dto.AuditResponse {
	return &dto.AuditResponse{
		ID:            entity.ID,
		Date:          *timeUtils.FormatResponseTime(&entity.Date),
		EntityName:    entity.EntityName,
		EntityId:      entity.EntityId,
		OldValue:      entity.OldValue,
		NewValue:      entity.NewValue,
		OperationType: entity.OperationType,
		Origin:        entity.Origin,
		UserID:        entity.UserID,
		User:          *entity.User.ToResponseModel(),
		PermissionID:  entity.PermissionID,
		Permission:    *entity.Permission.ToResponseModel(),
	}
}
