package permission

import (
	"context"

	"github.com/ericmarcelinotju/gram/dto"
)

// Service defines permission service behavior.
type Service interface {
	Create(ctx context.Context, payload *dto.PostPermissionDto) (*dto.PermissionDto, error)
	Read(ctx context.Context, filter *dto.GetPermissionDto) ([]dto.PermissionDto, int64, error)
	ReadById(context.Context, string) (*dto.PermissionDto, error)
	Update(context.Context, *dto.PutPermissionDto) (*dto.PermissionDto, error)
	DeleteById(context.Context, string) error
}

type service struct {
	repo Repository
}

// NewService creates a new service struct
func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (svc *service) Create(ctx context.Context, payload *dto.PostPermissionDto) (res *dto.PermissionDto, err error) {
	res = &dto.PermissionDto{
		Method:      payload.Method,
		Module:      payload.Module,
		Description: payload.Description,
	}
	err = svc.repo.Insert(ctx, res)
	return
}

func (svc *service) Read(ctx context.Context, payload *dto.GetPermissionDto) ([]dto.PermissionDto, int64, error) {
	return svc.repo.Select(
		ctx,
		&dto.PermissionDto{
			Method: *payload.Method,
			Module: *payload.Module,
		},
		payload.PaginationDto,
		payload.SortDto,
	)
}

func (svc *service) ReadById(ctx context.Context, id string) (*dto.PermissionDto, error) {
	return svc.repo.SelectById(ctx, id)
}

func (svc *service) Update(ctx context.Context, payload *dto.PutPermissionDto) (res *dto.PermissionDto, err error) {
	res = &dto.PermissionDto{
		Id:          payload.Id,
		Method:      payload.Method,
		Module:      payload.Module,
		Description: payload.Description,
	}
	err = svc.repo.Update(ctx, res)
	return
}

func (svc *service) DeleteById(ctx context.Context, id string) error {
	return svc.repo.Delete(ctx, &dto.PermissionDto{Id: id})
}
