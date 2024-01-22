package model

import (
	"github.com/ericmarcelinotju/gram/dto"
	uuid "github.com/satori/go.uuid"
)

// PermissionEntity struct defines the database model for a permission.
type PermissionEntity struct {
	Model
	Method      string
	Module      string
	Description string
}

func (PermissionEntity) PermissionEntity() string {
	return "permissions"
}

func NewPermissionEntity(dto *dto.PermissionDto) *PermissionEntity {
	id, _ := uuid.FromString(dto.Id)

	return &PermissionEntity{
		Model: Model{
			Id:        id,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
		Method:      dto.Method,
		Module:      dto.Module,
		Description: dto.Description,
	}
}

func (entity *PermissionEntity) ToDto() *dto.PermissionDto {
	return &dto.PermissionDto{
		Id:          entity.Id.String(),
		Method:      entity.Method,
		Module:      entity.Module,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
