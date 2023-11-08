package log

import (
	"context"
	"errors"

	pkgErr "github.com/pkg/errors"

	"github.com/ericmarcelinotju/gram/data/schema"
	domainErrors "github.com/ericmarcelinotju/gram/domain/errors"
	"github.com/ericmarcelinotju/gram/domain/model"

	"gorm.io/gorm"
)

const (
	insertError = "Error in inserting new log"
	updateError = "Error in updating log"
	deleteError = "Error in deleting log"
	selectError = "Error in selecting logs in the database"
)

type Store struct {
	db *gorm.DB
}

// New creates a new Store struct
func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) InsertLog(ctx context.Context, payload *model.Log) error {
	entity := schema.NewLogSchema(payload)

	if err := s.db.WithContext(ctx).Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.DatabaseError)
		return appErr
	}

	payload.ID = entity.ID.String()

	return nil
}

func (s *Store) SelectLog(ctx context.Context, filter *model.Log) ([]model.Log, error) {
	var entities []schema.Log

	query := s.db.WithContext(ctx).
		Model(&schema.Log{})

	if filter != nil {
		query.Where(schema.NewLogSchema(filter))
	}
	query.
		Order("created_at DESC").
		Find(&entities)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), domainErrors.NotFoundError)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.DatabaseError)
		return nil, appErr
	}

	var results = make([]model.Log, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDomainModel()
	}

	return results, nil
}

func (s *Store) SelectLogByID(ctx context.Context, id string) (*model.Log, error) {
	var log schema.Log

	if err := s.db.WithContext(ctx).
		Model(&schema.Log{}).
		First(&log, "id = ?", id).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.DatabaseError)
		return nil, appErr
	}

	return log.ToDomainModel(), nil
}

func (s *Store) DeleteLog(ctx context.Context, payload *model.Log) error {
	entity := schema.NewLogSchema(payload)

	if err := s.db.WithContext(ctx).Delete(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, deleteError), domainErrors.DatabaseError)
		return appErr
	}

	return nil
}
