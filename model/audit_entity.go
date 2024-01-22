package model

import (
	uuid "github.com/satori/go.uuid"
	"time"

	"github.com/ericmarcelinotju/gram/dto"
)

// AuditEntity struct defines the database model for a audit.
type AuditEntity struct {
	Id            uuid.UUID `gorm:"type:string"`
	Date          time.Time
	EntityName    string
	EntityId      string
	OldValue      string
	NewValue      string
	OperationType string
	Origin        string
	UserId        uuid.UUID
	User          UserEntity `gorm:"foreignKey:UserId"`
	PermissionId  uuid.UUID
	Permission    PermissionEntity `gorm:"foreignKey:PermissionId"`
}

func (AuditEntity) TableName() string {
	return "audits"
}

func NewAuditEntity(dto *dto.AuditDto) *AuditEntity {
	id, _ := uuid.FromString(dto.ID)
	userId, _ := uuid.FromString(dto.UserID)
	permissionID, _ := uuid.FromString(dto.PermissionID)

	return &AuditEntity{
		Id:            id,
		Date:          dto.Date,
		EntityName:    dto.EntityName,
		EntityId:      dto.EntityId,
		OldValue:      dto.OldValue,
		NewValue:      dto.NewValue,
		OperationType: dto.OperationType,
		Origin:        dto.Origin,
		UserId:        userId,
		PermissionId:  permissionID,
	}
}

func (entity *AuditEntity) ToDto() *dto.AuditDto {

	return &dto.AuditDto{
		ID:            entity.Id.String(),
		Date:          entity.Date,
		EntityName:    entity.EntityName,
		EntityId:      entity.EntityId,
		OldValue:      entity.OldValue,
		NewValue:      entity.NewValue,
		OperationType: entity.OperationType,
		Origin:        entity.Origin,
		UserID:        entity.UserId.String(),
		User:          *entity.User.ToDto(),
		PermissionID:  entity.PermissionId.String(),
		Permission:    *entity.Permission.ToDto(),
	}
}
