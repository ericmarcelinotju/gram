package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkgErr "github.com/pkg/errors"

	domainErrors "github.com/ericmarcelinotju/gram/domain/errors"
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/library/email"
	"github.com/ericmarcelinotju/gram/model"
	"github.com/ericmarcelinotju/gram/repository/cache"
	"github.com/ericmarcelinotju/gram/repository/notifier"
	"github.com/ericmarcelinotju/gram/utils/crypt"
)

const (
	loginError  = "error in attempting login"
	forgotError = "error in forgot password process"
)

// Repository provides an abstraction on top of the building data source
type Repository interface {
	Login(ctx context.Context, username string, password string, isRememberMe bool) (*dto.UserDto, string, error)
	Logout(context.Context, string) error
	ReadUserByToken(context.Context, string) (*dto.UserDto, error)
	ForgotPassword(context.Context, *dto.UserDto) error
}

type repository struct {
	db       *gorm.DB
	cache    cache.Cache
	notifier notifier.Notifier
}

// New creates a new Store struct
func NewRepository(
	db *gorm.DB,
	cache cache.Cache,
	notifier notifier.Notifier,
) *repository {
	return &repository{
		db:       db,
		cache:    cache,
		notifier: notifier,
	}
}

func (s *repository) Login(ctx context.Context, username string, password string, isRememberMe bool) (*dto.UserDto, string, error) {
	var err error
	var result model.UserEntity

	query := s.db.
		Preload("Role").
		Preload("Role.Permissions").
		First(&result, "username = ?", username)
	if err = query.Error; err != nil {
		err = domainErrors.NewAppError(pkgErr.Wrap(err, loginError), domainErrors.NotAuthorized)
		return nil, "", err
	}
	if !crypt.CompareHash(result.Password, password) {
		err = domainErrors.NewAppError(pkgErr.Wrap(err, loginError), domainErrors.NotAuthorized)
		return nil, "", err
	}

	token := uuid.New().String()
	expiry := time.Hour * 24
	if isRememberMe {
		expiry = time.Hour * 730
	}
	err = s.cache.Set(ctx, token, result.ToDto(), expiry)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	result.LastLogin = &now
	if err = s.db.WithContext(ctx).Model(&result).Updates(result).Error; err != nil {
		err = domainErrors.NewAppError(pkgErr.Wrap(err, loginError), domainErrors.DatabaseError)
		return nil, "", err
	}

	ginCtx, ok := ctx.(*gin.Context)
	if ok {
		http.SetCookie(ginCtx.Writer, &http.Cookie{
			Name:    "auth",
			Value:   token,
			Path:    "/",
			Expires: time.Now().Add(expiry),
			MaxAge:  int(expiry.Seconds()),
		})
	}

	return result.ToDto(), token, nil
}

func (s *repository) Logout(ctx context.Context, token string) error {
	err := s.cache.Get(ctx, token, nil)
	if err != nil {
		return domainErrors.NewAppError(errors.New("token not registered"), domainErrors.NotAuthorized)
	}
	err = s.cache.Del(ctx, token)
	if err != nil {
		return domainErrors.NewAppError(errors.New("failed to delete token"), domainErrors.NotAuthorized)
	}

	ginCtx, ok := ctx.(*gin.Context)
	if ok {
		http.SetCookie(ginCtx.Writer, &http.Cookie{
			Name:    "auth",
			Value:   "",
			Path:    "/",
			Expires: time.Now(),
			MaxAge:  -1,
		})
	}

	return nil
}

func (s *repository) ReadUserByToken(ctx context.Context, token string) (user *dto.UserDto, err error) {
	err = s.cache.Get(ctx, token, &user)
	if err != nil {
		err = domainErrors.NewAppError(errors.New("token not registered"), domainErrors.NotAuthorized)
	}
	return
}

func (s *repository) ForgotPassword(ctx context.Context, payload *dto.UserDto) error {
	token := uuid.New().String()
	err := s.cache.Set(ctx, token, payload, time.Minute*30)
	if err != nil {
		return domainErrors.NewAppError(pkgErr.Wrap(err, forgotError), domainErrors.CacheError)
	}

	err = s.notifier.Notify(
		"Password Reset",
		email.EmailContent{Data: os.Getenv("FRONTEND_URL") + "#/forgot-password?fpkey=" + token},
		payload,
	)
	if err != nil {
		return domainErrors.NewAppError(pkgErr.Wrap(err, forgotError), domainErrors.RepositoryError)
	}
	return nil
}
