package media

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gitlab.com/firelogik/helios/config"
	"gitlab.com/firelogik/helios/data/storage"
	domainErrors "gitlab.com/firelogik/helios/domain/errors"
)

type Store struct {
	storage storage.Storage
}

// New creates a new Store struct
func New(storage storage.Storage) *Store {
	return &Store{storage: storage}
}

func (s *Store) Upload(ctx context.Context, file *multipart.File) (string, error) {
	if file == nil {
		appErr := domainErrors.NewAppError(errors.New("uploaded file empty"), domainErrors.ValidationError)
		return "", appErr
	}
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), uuid.New().String())
	if err := s.storage.Upload(*file, filename); err != nil {
		appErr := domainErrors.NewAppError(fmt.Errorf("upload file failed: %w", err), domainErrors.RepositoryError)
		return "", appErr
	}
	return config.Get().MediaStorage.URL + filename, nil
}
