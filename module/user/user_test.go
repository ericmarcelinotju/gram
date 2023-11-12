package user

import (
	"context"
	"testing"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/repository/database"
	"github.com/ericmarcelinotju/gram/repository/storage"
	"github.com/go-playground/assert/v2"
)

func setupService() (context.Context, Service) {
	// TODO :: Use sqlite and populate data

	// get configuration stucts via .env file
	configuration := config.NewConfig()

	// establish DB connection
	db, _ := database.Connect(configuration.Database)
	var fileStorage storage.Storage
	if configuration.MediaStorage != nil {
		// initialize File Manager
		fileStorage, _ = storage.InitFile(configuration.MediaStorage)
	}

	userRepo := NewRepository(db, fileStorage)
	return context.Background(), NewService(userRepo)
}

func TestReadUserHandler(t *testing.T) {
	ctx, svc := setupService()

	res, total, err := svc.Read(ctx, nil)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, total, 0)
	assert.Equal(t, len(res), total)
}

func TestReadWithPaginationUserHandler(t *testing.T) {
	ctx, svc := setupService()

	currentPage := 1
	totalPerPage := 10

	res, total, err := svc.Read(ctx, &dto.GetUserDto{
		PaginationDto: &dto.PaginationDto{
			Limit: &totalPerPage,
			Page:  &currentPage,
		},
	})

	assert.NotEqual(t, err, nil)
	assert.Equal(t, total, 0)
	assert.Equal(t, len(res), totalPerPage)
}

func TestReadByIdUserHandler(t *testing.T) {
	ctx, svc := setupService()

	id := "asdasdasd"

	res, err := svc.ReadById(ctx, id)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, res.Id, id)
}

func TestReadByUsernameUserHandler(t *testing.T) {
	ctx, svc := setupService()

	username := "asdasdasd"

	res, err := svc.ReadByUsername(ctx, username)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, res.Username, username)
}

func TestCreateUserHandler(t *testing.T) {
	ctx, svc := setupService()

	payload := dto.PostUserDto{
		Username: "testing",
		Email:    "testing@gmail.com",
		Password: "password",
	}

	res, err := svc.Create(ctx, &payload)

	assert.NotEqual(t, err, nil)

	check, err := svc.ReadById(ctx, res.Id)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, res.Id, check.Id)
}

func TestUpdateUserHandler(t *testing.T) {
	ctx, svc := setupService()

	payload := dto.PutUserDto{
		Id:       "asdasdasd",
		Username: "testing-updated",
		Email:    "testing@gmail.com",
	}

	res, err := svc.Update(ctx, &payload)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, res.Id, payload.Id)
	assert.Equal(t, res.Username, payload.Username)
}

func TestDeleteUserHandler(t *testing.T) {
	ctx, svc := setupService()

	id := "aasdasdasd"

	err := svc.DeleteById(ctx, id)

	assert.NotEqual(t, err, nil)

	_, err = svc.ReadById(ctx, id)

	assert.Equal(t, err, nil)
}
