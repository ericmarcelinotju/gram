package permission

import (
	"context"
	"errors"
	customErrors "github.com/ericmarcelinotju/gram/errors"
	pkgErr "github.com/pkg/errors"
	"testing"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/plugins/database"
	"github.com/go-playground/assert/v2"
)

func setupService() (context.Context, Service) {
	// TODO :: Use sqlite and populate data

	// get configuration stucts via .env file
	configuration := config.NewConfig(".env.test")

	// establish DB connection
	db, _ := database.Connect(configuration.Database)

	repo := NewRepository(db)
	return context.Background(), NewService(repo)
}

func TestReadUserHandler(t *testing.T) {
	ctx, svc := setupService()

	res, total, err := svc.Read(ctx, &dto.GetPermissionDto{
		Method:        nil,
		Module:        nil,
		PaginationDto: nil,
		SortDto:       nil,
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(res), int(total))
}

func TestReadWithPaginationUserHandler(t *testing.T) {
	ctx, svc := setupService()

	totalPerPage := 10
	currentPage := 1

	res, total, err := svc.Read(ctx, &dto.GetPermissionDto{
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

func TestCreateUserHandler(t *testing.T) {
	ctx, svc := setupService()

	payload := dto.PostPermissionDto{
		Method:      "testing",
		Module:      "testing",
		Description: "testing@gmail.com",
	}

	res, err := svc.Create(ctx, &payload)

	assert.Equal(t, err, nil)

	check, err := svc.ReadById(ctx, res.Id)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Id, check.Id)
}

func TestUpdateUserHandler(t *testing.T) {
	ctx, svc := setupService()
	findQuery := "testing"
	permissions, total, err := svc.Read(ctx, &dto.GetPermissionDto{
		Method:        &findQuery,
		Module:        &findQuery,
		PaginationDto: nil,
		SortDto:       nil,
	})
	assert.NotEqual(t, total, 0)

	payload := dto.PutPermissionDto{
		Id:          permissions[0].Id,
		Method:      "testing-updated",
		Module:      "testing-updated",
		Description: "testing@gmail.com",
	}

	res, err := svc.Update(ctx, &payload)

	assert.Equal(t, err, nil)

	check, err := svc.ReadById(ctx, res.Id)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Id, check.Id)
	assert.Equal(t, res.Method, check.Method)
}

func TestReadByIdUserHandler(t *testing.T) {
	ctx, svc := setupService()

	permissions, total, err := svc.Read(ctx, &dto.GetPermissionDto{
		Method:        nil,
		Module:        nil,
		PaginationDto: nil,
		SortDto:       nil,
	})
	assert.NotEqual(t, total, 0)
	res, err := svc.ReadById(ctx, permissions[0].Id)

	assert.Equal(t, err, nil)
	assert.NotEqual(t, res, nil)
	assert.Equal(t, res.Id, permissions[0].Id)
}

func TestDeleteUserHandler(t *testing.T) {
	ctx, svc := setupService()

	findQuery := "testing-updated"
	permissions, total, err := svc.Read(ctx, &dto.GetPermissionDto{
		Method:        &findQuery,
		Module:        &findQuery,
		PaginationDto: nil,
		SortDto:       nil,
	})
	assert.NotEqual(t, total, 0)
	res, err := svc.ReadById(ctx, permissions[0].Id)

	err = svc.DeleteById(ctx, res.Id)

	assert.Equal(t, err, nil)
	permission, err := svc.ReadById(ctx, res.Id)

	appErr := customErrors.NewAppError(pkgErr.Wrap(errors.New("record not found"), selectError), customErrors.DatabaseError)
	assert.Equal(t, err.Error(), appErr.Error())
	assert.Equal(t, permission, nil)
}
