package permission

import (
	"context"

	"gitlab.com/firelogik/helios/domain/model"
)

// Repository provides an abstraction on top of the permission data source
type Repository interface {
	InsertPermission(context.Context, *model.Permission) error
	UpdatePermission(context.Context, *model.Permission) error
	SelectPermission(context.Context, *model.Permission) ([]model.Permission, int64, error)
	SelectPermissionByID(context.Context, string) (*model.Permission, error)
	DeletePermission(context.Context, *model.Permission) error
}
