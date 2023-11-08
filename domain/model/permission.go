package model

import (
	"time"

	dto "gitlab.com/firelogik/helios/router/http/dto/responses"
	timeUtils "gitlab.com/firelogik/helios/utils/time"
)

// Permission struct defines the database model for a permission.
type Permission struct {
	ID          string
	Method      string
	Module      string
	Description string

	Pagination Pagination
	Sort       Sort

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (entity *Permission) ToResponseModel() *dto.PermissionResponse {
	return &dto.PermissionResponse{
		ID:          entity.ID,
		Method:      entity.Method,
		Module:      entity.Module,
		Description: entity.Description,
		CreatedAt:   *timeUtils.FormatResponseTime(&entity.CreatedAt),
		UpdatedAt:   *timeUtils.FormatResponseTime(&entity.UpdatedAt),
	}
}
