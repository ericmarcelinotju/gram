package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ericmarcelinotju/gram/domain/model"
)

// Role struct defines the database model for a role.
type Role struct {
	Model
	Name        string `gorm:"unique"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

// RolePermissions struct defines the database model for a role permission.
type RolePermission struct {
	RoleID       uuid.UUID  `gorm:"primaryKey;type:uuid"`
	Role         Role       `gorm:"foreignKey:RoleID"`
	PermissionID uuid.UUID  `gorm:"primaryKey;type:uuid"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func NewRoleSchema(entity *model.Role) *Role {
	var permissions = make([]Permission, len(entity.Permissions))
	for i, permission := range entity.Permissions {
		permissions[i] = *NewPermissionSchema(&permission)
	}

	id, _ := uuid.Parse(entity.ID)

	return &Role{
		Model: Model{
			ID:        id,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Name:        entity.Name,
		Description: entity.Description,
		Permissions: permissions,
	}
}

func (entity *Role) ToDomainModel() *model.Role {
	var permissions = make([]model.Permission, len(entity.Permissions))
	for i, permission := range entity.Permissions {
		permissions[i] = *permission.ToDomainModel()
	}

	return &model.Role{
		ID:          entity.ID.String(),
		Name:        entity.Name,
		Description: entity.Description,
		Permissions: permissions,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
