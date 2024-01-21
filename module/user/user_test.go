package user

import (
	"context"
	"errors"
	customErrors "github.com/ericmarcelinotju/gram/errors"
	pkgErr "github.com/pkg/errors"
	"testing"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/plugins/database"
	"github.com/ericmarcelinotju/gram/plugins/storage"
	"github.com/go-playground/assert/v2"
)

func setupService() (context.Context, Service) {
	// TODO :: Use sqlite and populate data

	// get configuration stucts via .env file
	configuration := config.NewConfig(".env.test")

	// establish DB connection
	db, _ := database.Connect(configuration.Database)
	var fileStorage storage.Storage
	if configuration.MediaStorage != nil {
		// initialize File Manager
		fileStorage, _ = storage.NewFileStorage(configuration.MediaStorage)
	}

	userRepo := NewRepository(db, fileStorage, nil)
	return context.Background(), NewService(userRepo)
}

func TestReadUserHandler(t *testing.T) {
	ctx, svc := setupService()

	res, total, err := svc.Read(ctx, &dto.GetUserDto{
		Name:          nil,
		Email:         nil,
		RoleId:        nil,
		PaginationDto: nil,
		SortDto:       nil,
	})

	assert.Equal(t, err, nil)
	assert.Equal(t, len(res), int(total))
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

	assert.Equal(t, err, nil)
	if int(total) > totalPerPage {
		assert.Equal(t, len(res), totalPerPage)
	}
}

func TestReadByIdUserHandler(t *testing.T) {
	ctx, svc := setupService()

	res, total, err := svc.Read(ctx, &dto.GetUserDto{
		Name:  nil,
		Email: nil,
	})
	assert.NotEqual(t, total, 0)

	user, err := svc.ReadById(ctx, res[0].Id)

	assert.Equal(t, err, nil)
	assert.Equal(t, user.Id, res[0].Id)
}

func TestReadByUsernameUserHandler(t *testing.T) {
	ctx, svc := setupService()

	res, total, err := svc.Read(ctx, &dto.GetUserDto{
		Name:  nil,
		Email: nil,
	})
	assert.NotEqual(t, total, 0)

	user, err := svc.ReadByUsername(ctx, res[0].Name)

	assert.Equal(t, err, nil)
	assert.Equal(t, res[0].Name, user.Name)
}

func TestCreateUserHandler(t *testing.T) {
	ctx, svc := setupService()

	payload := dto.PostUserDto{
		Name:     "testing",
		Email:    "testing@gmail.com",
		Password: "password",
	}

	res, err := svc.Create(ctx, &payload)

	assert.Equal(t, err, nil)

	check, err := svc.ReadById(ctx, res.Id)

	assert.Equal(t, err, nil)

	assert.Equal(t, res.Id, check.Id)

}

func TestUpdateUserHandler(t *testing.T) {
	ctx, svc := setupService()
	name := "testing"
	email := "testing@gmail.com"
	users, total, err := svc.Read(ctx, &dto.GetUserDto{
		Name:  &name,
		Email: &email,
	})
	assert.NotEqual(t, total, 0)

	payload := dto.PutUserDto{
		Id:    users[0].Id,
		Name:  "testing-updated",
		Email: "testing@gmail.com",
	}

	res, err := svc.Update(ctx, &payload)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Id, payload.Id)
	assert.Equal(t, res.Name, payload.Name)
}

func TestDeleteUserHandler(t *testing.T) {
	ctx, svc := setupService()

	name := "testing-updated"
	email := "testing@gmail.com"
	users, total, err := svc.Read(ctx, &dto.GetUserDto{
		Name:  &name,
		Email: &email,
	})
	assert.NotEqual(t, total, 0)

	err = svc.DeleteById(ctx, users[0].Id)

	assert.Equal(t, err, nil)

	_, err = svc.ReadById(ctx, users[0].Id)
	appErr := customErrors.NewAppError(pkgErr.Wrap(errors.New("record not found"), selectError), customErrors.DatabaseError)

	assert.NotEqual(t, err.Error(), appErr.Error())
}
