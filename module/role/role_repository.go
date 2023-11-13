package role

import (
	"context"
	"errors"

	pkgErr "github.com/pkg/errors"

	"github.com/ericmarcelinotju/gram/dto"
	customErrors "github.com/ericmarcelinotju/gram/errors"
	"github.com/ericmarcelinotju/gram/model"

	"gorm.io/gorm"
)

const (
	insertError = "Error in inserting new role"
	updateError = "Error in updating role"
	deleteError = "Error in deleting role"
	selectError = "Error in selecting roles in the database"
)

// Repository provides an abstraction on top of the role data source
type Repository interface {
	Insert(context.Context, *dto.RoleDto) error
	Update(context.Context, *dto.RoleDto) error
	Select(context.Context, *dto.RoleDto, *dto.PaginationDto, *dto.SortDto) ([]dto.RoleDto, int64, error)
	SelectById(context.Context, string) (*dto.RoleDto, error)
	Delete(context.Context, *dto.RoleDto) error
}

type repository struct {
	db *gorm.DB
}

// New creates a new Store struct
func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (s *repository) Insert(ctx context.Context, payload *dto.RoleDto) error {
	entity := model.NewRoleEntity(payload)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Permissions").Create(entity).Error; err != nil {
			appErr := customErrors.NewAppError(pkgErr.Wrap(err, insertError), customErrors.DatabaseError)
			return appErr
		}
		if err := tx.Model(entity).Association("Permissions").Append(entity.Permissions); err != nil {
			appErr := customErrors.NewAppError(pkgErr.Wrap(err, updateError), customErrors.DatabaseError)
			return appErr
		}
		payload.Id = entity.Id.String()
		return nil
	})
}

func (s *repository) Update(ctx context.Context, payload *dto.RoleDto) error {
	entity := model.NewRoleEntity(payload)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(entity).Omit("Permissions").Updates(entity).Error; err != nil {
			appErr := customErrors.NewAppError(pkgErr.Wrap(err, updateError), customErrors.DatabaseError)
			return appErr
		}

		if err := tx.Model(entity).Association("Permissions").Replace(entity.Permissions); err != nil {
			appErr := customErrors.NewAppError(pkgErr.Wrap(err, updateError), customErrors.DatabaseError)
			return appErr
		}

		return nil
	})
}

func (s *repository) Select(
	ctx context.Context,
	filter *dto.RoleDto,
	pagination *dto.PaginationDto,
	sort *dto.SortDto,
) ([]dto.RoleDto, int64, error) {

	var total int64
	var entities []model.RoleEntity

	query := s.db.WithContext(ctx).
		Model(&model.RoleEntity{}).
		Preload("Permissions")

	if filter != nil {
		query.Where(model.NewRoleEntity(filter))
	}
	query.Count(&total)
	if pagination != nil {
		pagination.Apply(query)
	}
	if sort != nil {
		sort.Apply(query)
	}
	query.Find(&entities)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := customErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), customErrors.NotFoundError)
		return nil, total, appErr
	}

	if err := query.Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.DatabaseError)
		return nil, total, appErr
	}

	var results = make([]dto.RoleDto, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDto()
	}

	return results, total, nil
}

func (s *repository) SelectById(ctx context.Context, id string) (*dto.RoleDto, error) {
	var entity model.RoleEntity

	if err := s.db.WithContext(ctx).
		Model(&model.RoleEntity{}).
		Preload("Permissions").
		First(&entity, "id = ?", id).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.DatabaseError)
		return nil, appErr
	}

	return entity.ToDto(), nil
}

func (s *repository) Delete(ctx context.Context, payload *dto.RoleDto) error {
	entity := model.NewRoleEntity(payload)

	if err := s.db.WithContext(ctx).Delete(entity).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, deleteError), customErrors.DatabaseError)
		return appErr
	}

	return nil
}
