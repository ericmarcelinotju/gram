package user

import (
	"context"
	"fmt"
	"time"

	domainErrors "github.com/ericmarcelinotju/gram/domain/errors"
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/utils/crypt"
)

// Service defines user service behavior.
type Service interface {
	Create(context.Context, *dto.PostUserDto) (*dto.UserDto, error)
	Read(context.Context, *dto.GetUserDto) ([]dto.UserDto, int64, error)
	ReadById(context.Context, string) (*dto.UserDto, error)
	ReadByUsername(context.Context, string) (*dto.UserDto, error)
	Update(context.Context, *dto.PutUserDto) (*dto.UserDto, error)
	UpdatePassword(context.Context, *dto.ChangeUserPasswordDto) error

	DeleteById(context.Context, string) error
}

type service struct {
	repo Repository
}

// NewService creates a new service struct
func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (svc *service) Create(ctx context.Context, payload *dto.PostUserDto) (res *dto.UserDto, err error) {
	var avatar *string
	if payload.Avatar != nil {
		file, err := payload.Avatar.Open()
		if err != nil {
			return nil, err
		}

		filename := "user/" + fmt.Sprintf("%d", time.Now().Unix())
		if err := svc.repo.SaveAvatar(&file, filename); err != nil {
			return nil, err
		}

		avatar = &filename
	}

	res = &dto.UserDto{
		Username:   payload.Username,
		Lastname:   payload.Lastname,
		Department: payload.Department,
		Title:      payload.Title,
		Email:      payload.Email,
		Password:   payload.Password,
		Avatar:     avatar,
		RoleId:     payload.RoleId,
	}
	err = svc.repo.Insert(ctx, res)
	return
}

func (svc *service) Read(ctx context.Context, payload *dto.GetUserDto) ([]dto.UserDto, int64, error) {
	return svc.repo.Select(
		ctx,
		&dto.UserDto{
			Username: *payload.Username,
			Email:    *payload.Email,
			RoleId:   *payload.RoleId,
		},
		payload.PaginationDto,
		payload.SortDto,
	)
}

func (svc *service) ReadByUsername(ctx context.Context, username string) (*dto.UserDto, error) {
	return svc.repo.SelectByUsername(ctx, username)
}

func (svc *service) ReadById(ctx context.Context, id string) (*dto.UserDto, error) {
	return svc.repo.SelectById(ctx, id)
}

func (svc *service) Update(ctx context.Context, payload *dto.PutUserDto) (res *dto.UserDto, err error) {
	var avatar *string
	if payload.Avatar != nil {
		file, err := payload.Avatar.Open()
		if err != nil {
			return nil, err
		}

		filename := "user/" + fmt.Sprintf("%d", time.Now().Unix())
		if err := svc.repo.SaveAvatar(&file, filename); err != nil {
			return nil, err
		}

		avatar = &filename
	}

	res = &dto.UserDto{
		Username:   payload.Username,
		Lastname:   payload.Lastname,
		Department: payload.Department,
		Title:      payload.Title,
		Email:      payload.Email,
		Avatar:     avatar,
		RoleId:     payload.RoleId,
	}
	err = svc.repo.Insert(ctx, res)
	return
}

func (svc *service) UpdatePassword(ctx context.Context, payload *dto.ChangeUserPasswordDto) error {
	user, err := svc.repo.SelectById(ctx, payload.Id)
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.NotFoundError)
	}
	if crypt.CompareHash(user.Password, payload.OldPassword) {
		return domainErrors.NewAppError(err, domainErrors.NotAuthorized)
	}
	return svc.repo.UpdatePassword(ctx, payload.Id, payload.NewPassword)
}

func (svc *service) DeleteById(ctx context.Context, id string) error {
	payload := &dto.UserDto{Id: id}
	err := svc.repo.Delete(ctx, payload)
	if err != nil {
		return err
	}
	if payload.Avatar != nil {
		return svc.repo.RemoveAvatar(*payload.Avatar)
	}
	return nil
}
