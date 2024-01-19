package role

import (
	"context"
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

	assert.NotEqual(t, err, nil)
	assert.Equal(t, total, 0)
	assert.Equal(t, len(res), total)
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

	assert.NotEqual(t, err, nil)
	assert.Equal(t, total, 0)
	assert.Equal(t, len(res), totalPerPage)
}

func TestCreateUserHandler(t *testing.T) {
	ctx, svc := setupService()

	payload := dto.PostRoleDto{
		Name:        "testing",
		Description: "testing@gmail.com",
	}

	res, err := svc.Create(ctx, &payload)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, res.Name, payload.Name)
}

func TestUpdateUserHandler(t *testing.T) {
	ctx, svc := setupService()

	payload := dto.PutRoleDto{
		Id:          "asdasdasd",
		Name:        "testing-updated",
		Description: "testing@gmail.com",
	}

	res, err := svc.Update(ctx, &payload)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, res.Id, payload.Id)
	assert.Equal(t, res.Name, payload.Name)
}

func TestReadByIdUserHandler(t *testing.T) {
	ctx, svc := setupService()

	id := "asdasdasd"

	res, err := svc.ReadById(ctx, id)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Id, id)
}

func TestDeleteUserHandler(t *testing.T) {
	ctx, svc := setupService()

	id := "aasdasdasd"

	err := svc.DeleteById(ctx, id)

	assert.NotEqual(t, err, nil)

	_, err = svc.ReadById(ctx, id)

	assert.Equal(t, err, nil)
}
