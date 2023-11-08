package role

import (
	"context"

	"gitlab.com/firelogik/helios/domain/model"
)

// Repository provides an abstraction on top of the role data source
type Repository interface {
	InsertRole(context.Context, *model.Role) error
	UpdateRole(context.Context, *model.Role) error
	SelectRole(context.Context, *model.Role) ([]model.Role, int64, error)
	SelectRoleByID(context.Context, string) (*model.Role, error)
	DeleteRole(context.Context, *model.Role) error
}
