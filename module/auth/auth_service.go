package auth

import (
	"context"

	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/module/user"
)

type Service interface {
	Login(context.Context, *dto.LoginDto) (*dto.UserDto, string, error)
	Logout(ctx context.Context, token string) error
	ReadUserByToken(context.Context, string) (*dto.UserDto, error)

	// Generate forgot password token for reset password, send forgot password email
	ForgotPassword(context.Context, *dto.ForgotUserPasswordDto) error
	// Reset password with new password provided, validate it with forgot password token
	ResetPassword(context.Context, *dto.ResetUserPasswordDto) error
}

type service struct {
	repo     Repository
	userRepo user.Repository
}

// NewService creates a new service struct
func NewService(repo Repository, userRepo user.Repository) *service {
	return &service{repo: repo, userRepo: userRepo}
}

func (svc *service) Login(ctx context.Context, payload *dto.LoginDto) (*dto.UserDto, string, error) {
	return svc.repo.Login(ctx, payload.Username, payload.Password, payload.IsRememberMe)
}

func (svc *service) Logout(ctx context.Context, token string) error {
	return svc.repo.Logout(ctx, token)
}

func (svc *service) ReadUserByToken(ctx context.Context, token string) (*dto.UserDto, error) {
	return svc.repo.ReadUserByToken(ctx, token)
}

func (svc *service) ForgotPassword(ctx context.Context, payload *dto.ForgotUserPasswordDto) error {
	user, err := svc.userRepo.SelectByUsername(ctx, payload.Username)
	if err != nil {
		return err
	}

	if err = svc.repo.ForgotPassword(ctx, user); err != nil {
		return err
	}
	return nil
}

func (svc *service) ResetPassword(ctx context.Context, payload *dto.ResetUserPasswordDto) error {
	user, err := svc.ReadUserByToken(ctx, payload.ForgotToken)
	if err != nil {
		return err
	}
	return svc.userRepo.UpdatePassword(ctx, user.Id, payload.NewPassword)
}
