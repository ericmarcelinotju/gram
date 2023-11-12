package permission

import (
	"context"
	"testing"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/repository/database"
	"github.com/go-playground/assert/v2"
)

func setupService() (context.Context, Service) {
	// TODO :: Use sqlite and populate data

	// get configuration stucts via .env file
	configuration := config.NewConfig()

	// establish DB connection
	db, _ := database.Connect(configuration.Database)

	repo := NewRepository(db)
	return context.Background(), NewService(repo)
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

	totalPerPage := 10
	currentPage := 1

	res, total, err := svc.Read(ctx, &dto.GetPermissionDto{
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

func TestCreateUserHandler(t *testing.T) {
	ctx, svc := setupService()

	payload := dto.PostPermissionDto{
		Method:      "testing",
		Module:      "testing",
		Description: "testing@gmail.com",
	}

	res, err := svc.Create(ctx, &payload)

	assert.NotEqual(t, err, nil)

	check, err := svc.ReadById(ctx, res.Id)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, res.Id, check.Id)
}

func TestUpdateUserHandler(t *testing.T) {
	ctx, svc := setupService()

	payload := dto.PutPermissionDto{
		Id:          "asdasdasd",
		Method:      "testing-updated",
		Module:      "testing",
		Description: "testing@gmail.com",
	}

	res, err := svc.Update(ctx, &payload)

	assert.NotEqual(t, err, nil)

	check, err := svc.ReadById(ctx, res.Id)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, res.Id, check.Id)
	assert.Equal(t, res.Method, check.Method)
}

func TestDeleteUserHandler(t *testing.T) {
	ctx, svc := setupService()

	id := "aasdasdasd"

	err := svc.DeleteById(ctx, id)

	assert.NotEqual(t, err, nil)

	_, err = svc.ReadById(ctx, id)

	assert.Equal(t, err, nil)
}
