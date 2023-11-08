package log

import (
	"context"

	"github.com/ericmarcelinotju/gram/domain/model"
)

// Service defines log service behavior.
type Service interface {
	CreateLog(context.Context, *model.Log) error
	ReadLog(context.Context, *model.Log) ([]model.Log, error)
	ReadLogByID(context.Context, string) (*model.Log, error)
	DeleteLogByID(context.Context, string) error
}

type service struct {
	repo Repository
}

// NewService creates a new service struct
func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (svc *service) CreateLog(ctx context.Context, payload *model.Log) error {
	return svc.repo.InsertLog(ctx, payload)
}

func (svc *service) ReadLog(ctx context.Context, filter *model.Log) ([]model.Log, error) {
	return svc.repo.SelectLog(ctx, filter)
}

func (svc *service) ReadLogByID(ctx context.Context, id string) (*model.Log, error) {
	return svc.repo.SelectLogByID(ctx, id)
}

func (svc *service) DeleteLogByID(ctx context.Context, id string) error {
	return svc.repo.DeleteLog(ctx, &model.Log{ID: id})
}
