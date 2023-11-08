package schema

import (
	"time"

	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/google/uuid"
)

// Audit struct defines the database model for a audit.
type Audit struct {
	ID            uuid.UUID `gorm:"type:string"`
	Date          time.Time
	EntityName    string
	EntityId      string
	OldValue      string
	NewValue      string
	OperationType string
	Origin        string
	UserID        uuid.UUID
	User          User `gorm:"foreignKey:UserID"`
	PermissionID  uuid.UUID
	Permission    Permission `gorm:"foreignKey:PermissionID"`
}

func NewAuditSchema(entity *model.Audit) *Audit {
	id, _ := uuid.Parse(entity.ID)
	userId, _ := uuid.Parse(entity.UserID)
	permissionID, _ := uuid.Parse(entity.PermissionID)

	return &Audit{
		ID:            id,
		Date:          entity.Date,
		EntityName:    entity.EntityName,
		EntityId:      entity.EntityId,
		OldValue:      entity.OldValue,
		NewValue:      entity.NewValue,
		OperationType: entity.OperationType,
		Origin:        entity.Origin,
		UserID:        userId,
		PermissionID:  permissionID,
	}
}

func (entity *Audit) ToDomainModel() *model.Audit {
	return &model.Audit{
		ID:            entity.ID.String(),
		Date:          entity.Date,
		EntityName:    entity.EntityName,
		EntityId:      entity.EntityId,
		OldValue:      entity.OldValue,
		NewValue:      entity.NewValue,
		OperationType: entity.OperationType,
		Origin:        entity.Origin,
		UserID:        entity.UserID.String(),
		User:          *entity.User.ToDomainModel(),
		PermissionID:  entity.PermissionID.String(),
		Permission:    *entity.Permission.ToDomainModel(),
	}
}
