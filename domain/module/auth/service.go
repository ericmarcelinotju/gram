package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ericmarcelinotju/gram/constant/enums"
	domainErrors "github.com/ericmarcelinotju/gram/domain/errors"
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/ericmarcelinotju/gram/domain/module/log"
	"github.com/ericmarcelinotju/gram/domain/module/user"
	"github.com/ericmarcelinotju/gram/utils/crypt"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Login(*gin.Context, *model.User, bool) (string, error)
	Logout(*gin.Context, string) error
	ReadUserByToken(context.Context, string) (*model.User, error)

	// Generate forgot password token for reset password, send forgot password email
	ForgotPassword(context.Context, *model.User) error
	// Reset password with new password provided, validate it with forgot password token
	ResetPassword(context.Context, string, string) error
}

type service struct {
	repo     Repository
	userRepo user.Repository
	logRepo  log.Repository
}

// NewService creates a new service struct
func NewService(repo Repository, userRepo user.Repository, logRepo log.Repository) *service {
	return &service{repo: repo, userRepo: userRepo, logRepo: logRepo}
}

func (svc *service) Login(ctx *gin.Context, payload *model.User, isRememberMe bool) (string, error) {
	token, err := svc.repo.Login(ctx, payload, isRememberMe)
	if err != nil {
		return "", err
	}
	err = svc.logRepo.InsertLog(ctx, &model.Log{
		Title:   "User Logged In",
		Subject: fmt.Sprintf("User: %s logged in at %s", payload.Username, payload.LastLogin.Format("02 January 2006 15:04:05")),
		Content: fmt.Sprintf("User Information<br/>ID: %s<br/>Username: %s<br/>Email: %s<br/>Login Time: %s", payload.ID, payload.Username, payload.Email, payload.LastLogin.Format("02 January 2006 15:04:05")),
		Level:   enums.LogLevelInfo,
		Type:    enums.LogTypeEvent,
	})
	if err != nil {
		return "", err
	}

	return token, err
}

func (svc *service) Logout(ctx *gin.Context, token string) error {
	return svc.repo.Logout(ctx, token)
}

func (svc *service) ReadUserByToken(ctx context.Context, token string) (*model.User, error) {
	return svc.repo.ReadUserByToken(ctx, token)
}

func (svc *service) ForgotPassword(ctx context.Context, payload *model.User) error {
	result, _, err := svc.userRepo.SelectUser(ctx, payload)
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.NotFoundError)
	}
	if len(result) <= 0 {
		return domainErrors.NewAppError(errors.New("error in select user"), domainErrors.NotFoundError)
	}
	user := result[0]

	now := time.Now().Format("02 January 2006 15:04:05")
	err = svc.logRepo.InsertLog(ctx, &model.Log{
		Title:   "User Forgot Password",
		Subject: fmt.Sprintf("User: %s used forgot password at %s", payload.Username, now),
		Content: fmt.Sprintf("User Information<br/>ID: %s<br/>Username: %s<br/>Email: %s<br/>Time: %s", payload.ID, payload.Username, payload.Email, now),
		Level:   enums.LogLevelInfo,
		Type:    enums.LogTypeEvent,
	})
	if err != nil {
		return err
	}

	if err = svc.repo.ForgotPassword(ctx, &user); err != nil {
		return domainErrors.NewAppError(err, domainErrors.TokenGeneratorError)
	}
	return nil
}

func (svc *service) ResetPassword(ctx context.Context, newPassword, token string) error {
	user, err := svc.ReadUserByToken(ctx, token)
	if err != nil {
		return err
	}
	user.Password, err = crypt.Hash(newPassword)
	if err != nil {
		return err
	}

	now := time.Now().Format("02 January 2006 15:04:05")
	err = svc.logRepo.InsertLog(ctx, &model.Log{
		Title:   "User Reset Password",
		Subject: fmt.Sprintf("User: %s reset his/her password at %s", user.Username, now),
		Content: fmt.Sprintf("User Information<br/>ID: %s<br/>Username: %s<br/>Email: %s<br/>Time: %s", user.ID, user.Username, user.Email, now),
		Level:   enums.LogLevelInfo,
		Type:    enums.LogTypeEvent,
	})
	if err != nil {
		return err
	}

	return svc.userRepo.UpdateUser(ctx, user)
}
