package user

import (
	"context"
	"errors"
	"fmt"

	domainErrors "github.com/ericmarcelinotju/gram/domain/errors"
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/ericmarcelinotju/gram/domain/module/role"
	"github.com/ericmarcelinotju/gram/utils/crypt"
	"github.com/google/uuid"
)

// Service defines user service behavior.
type Service interface {
	CreateUser(context.Context, *model.User) error
	ReadUser(context.Context, *model.User) ([]model.User, int64, error)
	ReadUserByID(context.Context, string) (*model.User, error)
	ReadUserByUsername(context.Context, string) (*model.User, error)
	UpdateUser(context.Context, *model.User) error
	UpdateUserPassword(context.Context, *model.User, string) error

	DeleteUserByID(context.Context, string) error

	SubscribePushLog(context.Context, *model.User) error
	UnsubscribePushLog(context.Context, *model.User) error
}

type service struct {
	repo     Repository
	roleRepo role.Repository
}

// NewService creates a new service struct
func NewService(
	repo Repository,
	roleRepo role.Repository,
) *service {
	return &service{repo: repo, roleRepo: roleRepo}
}

func (svc *service) CreateUser(ctx context.Context, payload *model.User) error {
	var err error
	payload.Password, err = crypt.Hash(payload.Password)
	if err != nil {
		return err
	}
	payload.ID = uuid.New().String()
	if payload.AvatarFile != nil {
		if err := svc.repo.SaveAvatar(payload); err != nil {
			return err
		}
	}
	return svc.repo.InsertUser(ctx, payload)
}

func (svc *service) ReadUser(ctx context.Context, filter *model.User) ([]model.User, int64, error) {
	return svc.repo.SelectUser(ctx, filter)
}

func (svc *service) ReadUserByUsername(ctx context.Context, username string) (*model.User, error) {
	result, _, err := svc.repo.SelectUser(ctx, &model.User{Username: username})
	if err != nil {
		return nil, domainErrors.NewAppError(err, domainErrors.NotFoundError)
	}
	if len(result) == 0 {
		return nil, domainErrors.NewAppError(fmt.Errorf("no user with username : %s", username), domainErrors.NotFoundError)
	}
	return &result[0], nil
}

func (svc *service) ReadUserByID(ctx context.Context, id string) (*model.User, error) {
	return svc.repo.SelectUserByID(ctx, id)
}

func (svc *service) UpdateUser(ctx context.Context, payload *model.User) error {
	var err error
	if len(payload.Password) > 0 {
		payload.Password, err = crypt.Hash(payload.Password)
		if err != nil {
			return err
		}
	}
	if payload.AvatarFile != nil {
		if err := svc.repo.SaveAvatar(payload); err != nil {
			return err
		}
	}
	return svc.repo.UpdateUser(ctx, payload)
}

func (svc *service) UpdateUserPassword(ctx context.Context, payload *model.User, newPassword string) error {
	user, err := svc.repo.SelectUserByID(ctx, payload.ID)
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.NotFoundError)
	}
	if crypt.CompareHash(user.Password, payload.Password) {
		return domainErrors.NewAppError(err, domainErrors.NotAuthorized)
	}

	newPassword, err = crypt.Hash(newPassword)
	if err != nil {
		return err
	}
	user.Password = newPassword
	return svc.repo.UpdateUser(ctx, user)
}

func (svc *service) DeleteUserByID(ctx context.Context, id string) error {
	payload := &model.User{ID: id}
	err := svc.repo.DeleteUser(ctx, payload)
	if err != nil {
		return err
	}
	return svc.repo.RemoveAvatar(payload)
}

func (svc *service) SubscribePushLog(ctx context.Context, payload *model.User) error {
	result, _, err := svc.repo.SelectUser(ctx, &model.User{ID: payload.ID})
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.NotFoundError)
	}
	if len(result) == 0 {
		return domainErrors.NewAppError(errors.New("no user found"), domainErrors.NotFoundError)
	}
	user := result[0]
	user.NotificationToken = payload.NotificationToken
	if err = svc.repo.SubscribePushLog(&user); err != nil {
		return err
	}
	return svc.repo.UpdateUser(ctx, &user)
}

func (svc *service) UnsubscribePushLog(ctx context.Context, payload *model.User) error {
	result, _, err := svc.repo.SelectUser(ctx, &model.User{ID: payload.ID})
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.NotFoundError)
	}
	if len(result) == 0 {
		return domainErrors.NewAppError(errors.New("no user found"), domainErrors.NotFoundError)
	}
	user := result[0]
	return svc.repo.UnsubscribePushLog(&user)
}
