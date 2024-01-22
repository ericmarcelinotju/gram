package model

import (
	uuid "github.com/satori/go.uuid"
	"time"

	"github.com/ericmarcelinotju/gram/dto"
	"gorm.io/gorm"
)

// RoleEntity struct defines the database model for a role.
type RoleEntity struct {
	Model
	Name        string `gorm:"unique"`
	Description string
	Level       int
	Permissions []PermissionEntity `gorm:"many2many:role_permissions;"`
}

func (RoleEntity) TableName() string {
	return "roles"
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

func (RolePermissionEntity) TableName() string {
	return "role_permission"
}

func NewRoleEntity(entity *dto.RoleDto) *RoleEntity {
	var permissions = make([]PermissionEntity, len(entity.Permissions))
	for i, permission := range entity.Permissions {
		permissions[i] = *NewPermissionEntity(&permission)
	}

	id, _ := uuid.FromString(entity.Id)

	return &RoleEntity{
		Model: Model{
			Id:        id,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Name:        entity.Name,
		Description: entity.Description,
		Level:       entity.Level,
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
		Level:       entity.Level,
		Permissions: permissions,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
