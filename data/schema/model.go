package schema

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gitlab.com/firelogik/helios/domain/model"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:string"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		uuid := uuid.New().String()
		tx.Statement.SetColumn("ID", uuid)
	}

	ctx := tx.Statement.Context

	userCtx := ctx.Value("auth-user")

	user, ok := userCtx.(*model.User)
	if !ok || user == nil {
		return nil
	}

	permissionCtx := ctx.Value("permission")

	permission, ok := permissionCtx.(*model.Permission)
	if !ok || permission == nil {
		return nil
	}

	userId, _ := uuid.Parse(user.ID)
	permissionId, _ := uuid.Parse(permission.ID)

	newValue, _ := json.Marshal(tx.Statement.Dest)

	tx.Model(&Audit{}).Create(&Audit{
		ID:            uuid.New(),
		OperationType: "insert",
		EntityName:    tx.Statement.Table,
		EntityId:      m.ID.String(),
		NewValue:      string(newValue),
		Origin:        "",
		UserID:        userId,
		PermissionID:  permissionId,
		Date:          time.Now(),
	})

	return nil
}

func (m *Model) BeforeUpdate(tx *gorm.DB) error {
	ctx := tx.Statement.Context

	userCtx := ctx.Value("auth-user")

	user, ok := userCtx.(*model.User)
	if !ok || user == nil {
		return nil
	}

	permissionCtx := ctx.Value("permission")

	permission, ok := permissionCtx.(*model.Permission)
	if !ok || permission == nil {
		return nil
	}

	userId, _ := uuid.Parse(user.ID)
	permissionId, _ := uuid.Parse(permission.ID)

	result := map[string]interface{}{}
	tx.Model(tx.Statement.Model).Find(result, "id = ?", m.ID)

	oldValue, _ := json.Marshal(result)
	newValue, _ := json.Marshal(tx.Statement.Dest)

	tx.Model(&Audit{}).Create(&Audit{
		ID:            uuid.New(),
		OperationType: "update",
		EntityName:    tx.Statement.Table,
		EntityId:      m.ID.String(),
		OldValue:      string(oldValue),
		NewValue:      string(newValue),
		Origin:        "",
		UserID:        userId,
		PermissionID:  permissionId,
		Date:          time.Now(),
	})

	return nil
}

func (m *Model) BeforeDelete(tx *gorm.DB) error {
	ctx := tx.Statement.Context

	userCtx := ctx.Value("auth-user")

	user, ok := userCtx.(*model.User)
	if !ok || user == nil {
		return nil
	}

	permissionCtx := ctx.Value("permission")

	permission, ok := permissionCtx.(*model.Permission)
	if !ok || permission == nil {
		return nil
	}

	userId, _ := uuid.Parse(user.ID)
	permissionId, _ := uuid.Parse(permission.ID)

	tx.Model(&Audit{}).Create(&Audit{
		ID:            uuid.New(),
		OperationType: "delete",
		EntityName:    tx.Statement.Table,
		EntityId:      m.ID.String(),
		Origin:        "",
		UserID:        userId,
		PermissionID:  permissionId,
		Date:          time.Now(),
	})

	return nil
}
