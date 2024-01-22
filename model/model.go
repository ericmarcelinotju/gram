package model

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	Id        uuid.UUID `gorm:"type:uuid;default:(uuid_generate_v4())"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context

	userCtx := ctx.Value("auth-user")

	user, ok := userCtx.(*UserEntity)
	if !ok || user == nil {
		return nil
	}

	permissionCtx := ctx.Value("permission")

	permission, ok := permissionCtx.(*PermissionEntity)
	if !ok || permission == nil {
		return nil
	}

	userId := user.Model.Id
	permissionId := permission.Model.Id

	newValue, _ := json.Marshal(tx.Statement.Dest)

	tx.Model(&AuditEntity{}).Create(&AuditEntity{
		Id:            uuid.NewV4(),
		OperationType: "insert",
		EntityName:    tx.Statement.Table,
		EntityId:      m.Id.String(),
		NewValue:      string(newValue),
		Origin:        "",
		UserId:        userId,
		PermissionId:  permissionId,
		Date:          time.Now(),
	})

	return nil
}

func (m *Model) BeforeUpdate(tx *gorm.DB) error {
	ctx := tx.Statement.Context

	userCtx := ctx.Value("auth-user")

	user, ok := userCtx.(*UserEntity)
	if !ok || user == nil {
		return nil
	}

	permissionCtx := ctx.Value("permission")

	permission, ok := permissionCtx.(*PermissionEntity)
	if !ok || permission == nil {
		return nil
	}

	userId := user.Model.Id
	permissionId := permission.Model.Id

	result := map[string]interface{}{}
	tx.Model(tx.Statement.Model).Find(result, "id = ?", m.Id)

	oldValue, _ := json.Marshal(result)
	newValue, _ := json.Marshal(tx.Statement.Dest)

	tx.Model(&AuditEntity{}).Create(&AuditEntity{
		Id:            uuid.NewV4(),
		OperationType: "update",
		EntityName:    tx.Statement.Table,
		EntityId:      m.Id.String(),
		OldValue:      string(oldValue),
		NewValue:      string(newValue),
		Origin:        "",
		UserId:        userId,
		PermissionId:  permissionId,
		Date:          time.Now(),
	})

	return nil
}

func (m *Model) BeforeDelete(tx *gorm.DB) error {
	ctx := tx.Statement.Context

	userCtx := ctx.Value("auth-user")

	user, ok := userCtx.(*UserEntity)
	if !ok || user == nil {
		return nil
	}

	permissionCtx := ctx.Value("permission")

	permission, ok := permissionCtx.(*PermissionEntity)
	if !ok || permission == nil {
		return nil
	}

	tx.Model(&AuditEntity{}).Create(&AuditEntity{
		Id:            uuid.NewV4(),
		OperationType: "delete",
		EntityName:    tx.Statement.Table,
		EntityId:      m.Id.String(),
		Origin:        "",
		UserId:        user.Model.Id,
		PermissionId:  permission.Model.Id,
		Date:          time.Now(),
	})

	return nil
}
