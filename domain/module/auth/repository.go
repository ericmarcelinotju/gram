package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/model"
)

// Repository provides an abstraction on top of the building data source
type Repository interface {
	Login(*gin.Context, *model.User, bool) (string, error)
	Logout(*gin.Context, string) error
	ReadUserByToken(context.Context, string) (*model.User, error)
	ForgotPassword(context.Context, *model.User) error
}
