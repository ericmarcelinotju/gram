package permission

import (
	"context"

	"github.com/ericmarcelinotju/gram/domain/model"
)

// Service defines permission service behavior.
type Service interface {
	CreatePermission(context.Context, *model.Permission) error
	ReadPermission(context.Context, *model.Permission) ([]model.Permission, int64, error)
	ReadPermissionByID(context.Context, string) (*model.Permission, error)
	UpdatePermission(context.Context, *model.Permission) error
	DeletePermissionByID(context.Context, string) error
}

type service struct {
	repo Repository
}

// NewService creates a new service struct
func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (svc *service) CreatePermission(ctx context.Context, payload *model.Permission) error {
	return svc.repo.InsertPermission(ctx, payload)
}

func (svc *service) ReadPermission(ctx context.Context, filter *model.Permission) ([]model.Permission, int64, error) {
	return svc.repo.SelectPermission(ctx, filter)
}

func (svc *service) ReadPermissionByID(ctx context.Context, id string) (*model.Permission, error) {
	return svc.repo.SelectPermissionByID(ctx, id)
}

func (svc *service) UpdatePermission(ctx context.Context, payload *model.Permission) error {
	return svc.repo.UpdatePermission(ctx, payload)
}

func (svc *service) DeletePermissionByID(ctx context.Context, id string) error {
	return svc.repo.DeletePermission(ctx, &model.Permission{ID: id})
}
