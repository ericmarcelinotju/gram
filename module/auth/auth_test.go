package auth

import (
	"context"
	"testing"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/module/user"
	"github.com/ericmarcelinotju/gram/plugins/cache"
	"github.com/ericmarcelinotju/gram/plugins/database"
	"github.com/ericmarcelinotju/gram/plugins/notifier"
	"github.com/ericmarcelinotju/gram/plugins/storage"
	"github.com/go-playground/assert/v2"
)

func setupService() (context.Context, Service) {
	// TODO :: Use sqlite and populate data

	// get configuration stucts via .env file
	configuration := config.NewConfig()

	// establish DB connection
	db, _ := database.Connect(configuration.Database)

	// establish cache connection
	cache, _ := cache.ConnectRedis(configuration.Cache)

	// initialize Forgot Password Email
	var emailer notifier.Notifier

	var fileStorage storage.Storage
	if configuration.MediaStorage != nil {
		// initialize File Manager
		fileStorage, _ = storage.NewFileStorage(configuration.MediaStorage)
	}

	repo := NewRepository(db, cache, emailer)
	userRepo := user.NewRepository(db, fileStorage, nil)
	return context.Background(), NewService(repo, userRepo)
}

func TestLoginHandler(t *testing.T) {
	ctx, svc := setupService()

	username := ""
	password := ""

	user, token, err := svc.Login(ctx, &dto.LoginDto{
		Username: username,
		Password: password,
	})

	assert.NotEqual(t, err, nil)
	assert.NotEqual(t, len(token), 0)
	assert.Equal(t, user.Username, username)
}

func TestLogoutHandler(t *testing.T) {
	ctx, svc := setupService()

	username := ""
	password := ""

	user, token, err := svc.Login(ctx, &dto.LoginDto{
		Username: username,
		Password: password,
	})

	assert.NotEqual(t, err, nil)
	assert.NotEqual(t, len(token), 0)
	assert.Equal(t, user.Username, username)

	err = svc.Logout(ctx, token)

	assert.NotEqual(t, err, nil)
}
