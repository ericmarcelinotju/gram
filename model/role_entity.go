package model

import (
	"time"

	"github.com/ericmarcelinotju/gram/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleEntity struct defines the database model for a role.
type RoleEntity struct {
	Model
	Name        string `gorm:"unique"`
	Description string
	Permissions []PermissionEntity `gorm:"many2many:role_permissions;"`
}

// RolePermissions struct defines the database model for a role permission.
type RolePermissionEntity struct {
	RoleId       uuid.UUID        `gorm:"primaryKey;type:uuid"`
	Role         RoleEntity       `gorm:"foreignKey:RoleId"`
	PermissionId uuid.UUID        `gorm:"primaryKey;type:uuid"`
	Permission   PermissionEntity `gorm:"foreignKey:PermissionId"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func NewRoleEntity(entity *dto.RoleDto) *RoleEntity {
	var permissions = make([]PermissionEntity, len(entity.Permissions))
	for i, permission := range entity.Permissions {
		permissions[i] = *NewPermissionEntity(&permission)
	}

	id, _ := uuid.Parse(entity.Id)

	return &RoleEntity{
		Model: Model{
			Id:        id,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Name:        entity.Name,
		Description: entity.Description,
		Permissions: permissions,
	}
}

func (entity *RoleEntity) ToDto() *dto.RoleDto {
	var permissions = make([]dto.PermissionDto, len(entity.Permissions))
	for i, permission := range entity.Permissions {
		permissions[i] = *permission.ToDto()
	}

	return &dto.RoleDto{
		Id:          entity.Id.String(),
		Name:        entity.Name,
		Description: entity.Description,
		Permissions: permissions,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
