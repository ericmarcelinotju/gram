package role

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

	res, total, err := svc.Read(ctx, &dto.GetRoleDto{
		Name:          nil,
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

	res, total, err := svc.Read(ctx, &dto.GetRoleDto{
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

	payload := dto.PostRoleDto{
		Name:        "testing",
		Description: "testing@gmail.com",
	}

	res, err := svc.Create(ctx, &payload)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Name, payload.Name)
}

func TestUpdateUserHandler(t *testing.T) {
	ctx, svc := setupService()
	name := "testing"
	res, total, err := svc.Read(ctx, &dto.GetRoleDto{
		Name:          &name,
		PaginationDto: nil,
		SortDto:       nil,
	})
	assert.NotEqual(t, total, 0)

	payload := dto.PutRoleDto{
		Id:          res[0].Id,
		Name:        "testing-updated",
		Description: "testing@gmail.com",
	}

	role, err := svc.Update(ctx, &payload)

	assert.Equal(t, err, nil)
	assert.Equal(t, role.Id, payload.Id)
	assert.Equal(t, role.Name, payload.Name)
}

func TestReadByIdUserHandler(t *testing.T) {
	ctx, svc := setupService()
	name := "testing-updated"
	res, total, err := svc.Read(ctx, &dto.GetRoleDto{
		Name:          &name,
		PaginationDto: nil,
		SortDto:       nil,
	})
	assert.NotEqual(t, total, 0)

	role, err := svc.ReadById(ctx, res[0].Id)

	assert.Equal(t, err, nil)
	assert.Equal(t, role.Id, res[0].Id)
}

func TestDeleteUserHandler(t *testing.T) {
	ctx, svc := setupService()

	name := "testing-updated"
	res, total, err := svc.Read(ctx, &dto.GetRoleDto{
		Name:          &name,
		PaginationDto: nil,
		SortDto:       nil,
	})
	assert.NotEqual(t, total, 0)

	err = svc.DeleteById(ctx, res[0].Id)

	assert.Equal(t, err, nil)

	_, err = svc.ReadById(ctx, res[0].Id)
	appErr := customErrors.NewAppError(pkgErr.Wrap(errors.New("record not found"), selectError), customErrors.DatabaseError)

	assert.Equal(t, err.Error(), appErr.Error())
}
