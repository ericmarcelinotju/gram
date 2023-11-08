package role

import (
	"context"

	"gitlab.com/firelogik/helios/domain/model"
	"gitlab.com/firelogik/helios/domain/module/permission"
)

// Service defines role service behavior.
type Service interface {
	CreateRole(context.Context, *model.Role) error
	ReadRole(context.Context, *model.Role) ([]model.Role, int64, error)
	ReadRoleByID(context.Context, string) (*model.Role, error)
	UpdateRole(context.Context, *model.Role) error
	DeleteRoleByID(context.Context, string) error
}

type service struct {
	repo           Repository
	permissionRepo permission.Repository
}

// NewService creates a new service struct
func NewService(repo Repository, permissionRepo permission.Repository) *service {
	return &service{repo: repo, permissionRepo: permissionRepo}
}

func (svc *service) CreateRole(ctx context.Context, payload *model.Role) error {
	return svc.repo.InsertRole(ctx, payload)
}

func (svc *service) ReadRole(ctx context.Context, filter *model.Role) ([]model.Role, int64, error) {
	return svc.repo.SelectRole(ctx, filter)
}

func (svc *service) ReadRoleByID(ctx context.Context, id string) (*model.Role, error) {
	return svc.repo.SelectRoleByID(ctx, id)
}

func (svc *service) UpdateRole(ctx context.Context, payload *model.Role) error {
	return svc.repo.UpdateRole(ctx, payload)
}

func (svc *service) DeleteRoleByID(ctx context.Context, id string) error {
	return svc.repo.DeleteRole(ctx, &model.Role{ID: id})
}
