package schema

import (
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/google/uuid"
)

// Permission struct defines the database model for a permission.
type Permission struct {
	Model
	Method      string
	Module      string
	Description string
}

func NewPermissionSchema(entity *model.Permission) *Permission {
	id, _ := uuid.Parse(entity.ID)

	return &Permission{
		Model: Model{
			ID:        id,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Method:      entity.Method,
		Module:      entity.Module,
		Description: entity.Description,
	}
}

func (entity *Permission) ToDomainModel() *model.Permission {
	return &model.Permission{
		ID:          entity.ID.String(),
		Method:      entity.Method,
		Module:      entity.Module,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
