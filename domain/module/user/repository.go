package user

import (
	"context"

	"gitlab.com/firelogik/helios/domain/model"
)

// Repository provides an abstraction on top of the user data source
type Repository interface {
	InsertUser(context.Context, *model.User) error
	UpdateUser(context.Context, *model.User) error
	SelectUser(context.Context, *model.User) ([]model.User, int64, error)
	SelectUserByID(context.Context, string) (*model.User, error)
	DeleteUser(context.Context, *model.User) error

	SaveAvatar(*model.User) error
	RemoveAvatar(*model.User) error

	SubscribePushLog(*model.User) error
	UnsubscribePushLog(*model.User) error
}
