package role

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
	insertError = "Error in inserting new role"
	updateError = "Error in updating role"
	deleteError = "Error in deleting role"
	selectError = "Error in selecting roles in the database"
)

type Store struct {
	db *gorm.DB
}

// New creates a new Store struct
func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Migrate(ctx context.Context) error {
	return s.db.AutoMigrate(&schema.Role{})
}

func (s *Store) InsertRole(ctx context.Context, payload *model.Role) error {
	entity := schema.NewRoleSchema(payload)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.db.WithContext(ctx).Omit("Permissions").Create(entity).Error; err != nil {
			appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.DatabaseError)
			return appErr
		}
		payload.ID = entity.ID.String()

		if err := tx.Model(entity).Association("Permissions").Append(entity.Permissions); err != nil {
			appErr := domainErrors.NewAppError(pkgErr.Wrap(err, updateError), domainErrors.DatabaseError)
			return appErr
		}
		return nil
	})
}

func (s *Store) UpdateRole(ctx context.Context, payload *model.Role) error {
	entity := schema.NewRoleSchema(payload)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(entity).Omit("Permissions").Updates(entity).Error; err != nil {
			appErr := domainErrors.NewAppError(pkgErr.Wrap(err, updateError), domainErrors.DatabaseError)
			return appErr
		}

		if err := tx.Model(entity).Association("Permissions").Replace(entity.Permissions); err != nil {
			appErr := domainErrors.NewAppError(pkgErr.Wrap(err, updateError), domainErrors.DatabaseError)
			return appErr
		}

		return nil
	})
}

func (s *Store) SelectRole(ctx context.Context, filter *model.Role) ([]model.Role, int64, error) {

	var total int64
	var entities []schema.Role

	query := s.db.WithContext(ctx).
		Model(&schema.Role{}).
		Preload("Permissions")

	if filter != nil {
		query.Where(schema.NewRoleSchema(filter))
	}
	query.Count(&total)
	if filter != nil {
		filter.Pagination.Apply(query)
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

	var results = make([]model.Role, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDomainModel()
	}

	return results, total, nil
}

func (s *Store) SelectRoleByID(ctx context.Context, id string) (*model.Role, error) {
	var role schema.Role

	if err := s.db.WithContext(ctx).
		Model(&schema.Role{}).
		Preload("Permissions").
		First(&role, "id = ?", id).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.DatabaseError)
		return nil, appErr
	}

	return role.ToDomainModel(), nil
}

func (s *Store) DeleteRole(ctx context.Context, payload *model.Role) error {
	entity := schema.NewRoleSchema(payload)

	if err := s.db.WithContext(ctx).Delete(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, deleteError), domainErrors.DatabaseError)
		return appErr
	}

	return nil
}
