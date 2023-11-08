package model

import (
	"time"

	dto "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	timeUtils "github.com/ericmarcelinotju/gram/utils/time"
)

// Role struct defines the database model for a role.
type Role struct {
	ID          string
	Name        string
	Description string
	Permissions []Permission

	Pagination Pagination
	Sort       Sort

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (entity *Role) ToResponseModel() *dto.RoleResponse {
	var permissions []dto.PermissionResponse = make([]dto.PermissionResponse, len(entity.Permissions))
	for i, item := range entity.Permissions {
		permissions[i] = *item.ToResponseModel()
	}

	return &dto.RoleResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Permissions: permissions,
		CreatedAt:   *timeUtils.FormatResponseTime(&entity.CreatedAt),
		UpdatedAt:   *timeUtils.FormatResponseTime(&entity.UpdatedAt),
	}
}
