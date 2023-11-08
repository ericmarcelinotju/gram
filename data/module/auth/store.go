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

	"github.com/ericmarcelinotju/gram/data/cache"
	"github.com/ericmarcelinotju/gram/data/notifier"
	"github.com/ericmarcelinotju/gram/data/schema"
	domainErrors "github.com/ericmarcelinotju/gram/domain/errors"
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/ericmarcelinotju/gram/library/email"
	"github.com/ericmarcelinotju/gram/utils/crypt"
)

const (
	loginError = "error in attempting login"
)

type Store struct {
	db       *gorm.DB
	cache    cache.Cache
	notifier notifier.Notifier
}

// New creates a new Store struct
func New(
	db *gorm.DB,
	cache cache.Cache,
	notifier notifier.Notifier,
) *Store {
	return &Store{
		db:       db,
		cache:    cache,
		notifier: notifier,
	}
}

func (s *Store) Login(ctx *gin.Context, payload *model.User, isRememberMe bool) (token string, err error) {
	var result schema.User

	query := s.db.
		Preload("Role").
		Preload("Role.Permissions").
		First(&result, "username = ?", payload.Username)
	if err = query.Error; err != nil {
		err = domainErrors.NewAppError(pkgErr.Wrap(err, loginError), domainErrors.NotAuthorized)
		return
	}
	if !crypt.CompareHash(result.Password, payload.Password) {
		err = domainErrors.NewAppError(pkgErr.Wrap(err, loginError), domainErrors.NotAuthorized)
		return
	}

	token = uuid.New().String()
	expiry := time.Hour * 24
	if isRememberMe {
		expiry = time.Hour * 730
	}
	err = s.cache.Set(ctx, token, result.ToDomainModel(), expiry)
	if err != nil {
		return
	}

	now := time.Now()
	result.LastLogin = &now
	if err = s.db.WithContext(ctx).Model(&result).Updates(result).Error; err != nil {
		err = domainErrors.NewAppError(pkgErr.Wrap(err, loginError), domainErrors.DatabaseError)
		return
	}

	*payload = *result.ToDomainModel()

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "auth",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(expiry),
		MaxAge:  int(expiry.Seconds()),
	})

	return
}

func (s *Store) Logout(ctx *gin.Context, token string) error {
	err := s.cache.Get(ctx, token, nil)
	if err != nil {
		return domainErrors.NewAppError(errors.New("token not registered"), domainErrors.NotAuthorized)
	}
	err = s.cache.Del(ctx, token)
	if err != nil {
		return domainErrors.NewAppError(errors.New("failed to delete token"), domainErrors.NotAuthorized)
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "auth",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
		MaxAge:  -1,
	})

	return nil
}

func (s *Store) ReadUserByToken(ctx context.Context, token string) (user *model.User, err error) {
	err = s.cache.Get(ctx, token, &user)
	if err != nil {
		err = domainErrors.NewAppError(errors.New("token not registered"), domainErrors.NotAuthorized)
	}
	return
}

func (s *Store) ForgotPassword(ctx context.Context, payload *model.User) error {
	token := uuid.New().String()
	err := s.cache.Set(ctx, token, payload, time.Minute*30)
	if err != nil {
		return err
	}

	return s.notifier.Notify(
		"Password Reset",
		email.EmailContent{Data: os.Getenv("FRONTEND_URL") + "#/forgot-password?fpkey=" + token},
		payload,
	)
}
