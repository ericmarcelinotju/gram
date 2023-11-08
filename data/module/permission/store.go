package permission

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
	insertError = "Error in inserting new permission"
	updateError = "Error in updating permission"
	deleteError = "Error in deleting permission"
	selectError = "Error in selecting permissions in the database"
)

type Store struct {
	db *gorm.DB
}

// New creates a new Store struct
func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Migrate(ctx context.Context) error {
	return s.db.AutoMigrate(&schema.Permission{})
}

func (s *Store) InsertPermission(ctx context.Context, payload *model.Permission) error {
	entity := schema.NewPermissionSchema(payload)

	if err := s.db.WithContext(ctx).Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.DatabaseError)
		return appErr
	}

	payload.ID = entity.ID.String()

	return nil
}

func (s *Store) UpdatePermission(ctx context.Context, payload *model.Permission) error {
	entity := schema.NewPermissionSchema(payload)

	if err := s.db.WithContext(ctx).Model(entity).Updates(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, updateError), domainErrors.DatabaseError)
		return appErr
	}

	return nil
}

func (s *Store) SelectPermission(ctx context.Context, filter *model.Permission) ([]model.Permission, int64, error) {
	var total int64
	var entities []schema.Permission

	query := s.db.WithContext(ctx).
		Model(&schema.Permission{})

	if filter != nil {
		query.Where(schema.NewPermissionSchema(filter))
	}
	query.Count(&total)
	if filter != nil {
		query = filter.Pagination.Apply(query)
		filter.Sort.Apply(query)
	}
	query.Find(&entities)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), domainErrors.NotFoundError)
		return nil, total, appErr
	}

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.DatabaseError)
		return nil, total, appErr
	}

	var results = make([]model.Permission, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDomainModel()
	}

	return results, total, nil
}

func (s *Store) SelectPermissionByID(ctx context.Context, id string) (*model.Permission, error) {
	var permission schema.Permission

	if err := s.db.WithContext(ctx).
		Model(&schema.Permission{}).
		First(&permission, "id = ?", id).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.DatabaseError)
		return nil, appErr
	}

	return permission.ToDomainModel(), nil
}

func (s *Store) DeletePermission(ctx context.Context, payload *model.Permission) error {
	entity := schema.NewPermissionSchema(payload)

	if err := s.db.WithContext(ctx).Delete(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, deleteError), domainErrors.DatabaseError)
		return appErr
	}

	return nil
}
