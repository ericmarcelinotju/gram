package role

import (
	"context"

	"github.com/ericmarcelinotju/gram/dto"
)

// Service defines role service behavior.
type Service interface {
	Create(context.Context, *dto.PostRoleDto) (*dto.RoleDto, error)
	Read(context.Context, *dto.GetRoleDto) ([]dto.RoleDto, int64, error)
	ReadById(context.Context, string) (*dto.RoleDto, error)
	Update(context.Context, *dto.PutRoleDto) (*dto.RoleDto, error)
	DeleteById(context.Context, string) error
}

type service struct {
	repo Repository
}

// NewService creates a new service struct
func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (svc *service) Create(ctx context.Context, payload *dto.PostRoleDto) (res *dto.RoleDto, err error) {
	var permissions []dto.PermissionDto = make([]dto.PermissionDto, len(payload.Permissions))
	for i, item := range payload.Permissions {
		permissions[i] = dto.PermissionDto{Id: item.Id}
	}
	res = &dto.RoleDto{
		Name:        payload.Name,
		Description: payload.Description,
		Permissions: permissions,
	}
	err = svc.repo.Insert(ctx, res)
	return
}

func (svc *service) Read(ctx context.Context, payload *dto.GetRoleDto) ([]dto.RoleDto, int64, error) {
	filter := &dto.RoleDto{}

	if payload.Name != nil {
		filter.Name = *payload.Name
	}
	return svc.repo.Select(
		ctx,
		filter,
		payload.PaginationDto,
		payload.SortDto,
	)
}

func (svc *service) ReadById(ctx context.Context, id string) (*dto.RoleDto, error) {
	return svc.repo.SelectById(ctx, id)
}

func (svc *service) Update(ctx context.Context, payload *dto.PutRoleDto) (res *dto.RoleDto, err error) {
	var permissions []dto.PermissionDto = make([]dto.PermissionDto, len(payload.Permissions))
	for i, item := range payload.Permissions {
		permissions[i] = dto.PermissionDto{Id: item.Id}
	}
	res = &dto.RoleDto{
		Id:          payload.Id,
		Name:        payload.Name,
		Description: payload.Description,
		Permissions: permissions,
	}
	err = svc.repo.Update(ctx, res)
	return
}

func (svc *service) DeleteById(ctx context.Context, id string) error {
	return svc.repo.Delete(ctx, &dto.RoleDto{Id: id})
}
