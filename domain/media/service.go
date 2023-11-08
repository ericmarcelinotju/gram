package media

import (
	"context"
	"mime/multipart"
)

// Service defines building service behavior.
type Service interface {
	Upload(context.Context, *multipart.File) (string, error)
}

type service struct {
	repo Repository
}

// NewService creates a new service struct
func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (svc *service) Upload(ctx context.Context, file *multipart.File) (string, error) {
	return svc.repo.Upload(ctx, file)
}
