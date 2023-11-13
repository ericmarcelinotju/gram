package permission

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
	insertError = "Error in inserting new permission"
	updateError = "Error in updating permission"
	deleteError = "Error in deleting permission"
	selectError = "Error in selecting permissions in the database"
)

// Repository provides an abstraction on top of the permission data source
type Repository interface {
	Insert(context.Context, *dto.PermissionDto) error
	Update(context.Context, *dto.PermissionDto) error
	Select(context.Context, *dto.PermissionDto, *dto.PaginationDto, *dto.SortDto) ([]dto.PermissionDto, int64, error)
	SelectById(context.Context, string) (*dto.PermissionDto, error)
	Delete(context.Context, *dto.PermissionDto) error
}

type repository struct {
	db *gorm.DB
}

// New creates a new Store struct
func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (s *repository) Insert(ctx context.Context, payload *dto.PermissionDto) error {
	entity := model.NewPermissionEntity(payload)

	if err := s.db.WithContext(ctx).Create(entity).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, insertError), customErrors.DatabaseError)
		return appErr
	}

	payload.Id = entity.Id.String()

	return nil
}

func (s *repository) Update(ctx context.Context, payload *dto.PermissionDto) error {
	entity := model.NewPermissionEntity(payload)

	if err := s.db.WithContext(ctx).Model(entity).Updates(entity).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, updateError), customErrors.DatabaseError)
		return appErr
	}

	return nil
}

func (s *repository) Select(
	ctx context.Context,
	filter *dto.PermissionDto,
	pagination *dto.PaginationDto,
	sort *dto.SortDto,
) ([]dto.PermissionDto, int64, error) {
	var total int64
	var entities []model.PermissionEntity

	query := s.db.WithContext(ctx).
		Model(&model.PermissionEntity{})

	if filter != nil {
		query.Where(model.NewPermissionEntity(filter))
	}
	query.Count(&total)
	if pagination != nil {
		query = pagination.Apply(query)
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

	var results = make([]dto.PermissionDto, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDto()
	}

	return results, total, nil
}

func (s *repository) SelectById(ctx context.Context, id string) (*dto.PermissionDto, error) {
	var permission model.PermissionEntity

	if err := s.db.WithContext(ctx).
		Model(&model.PermissionEntity{}).
		First(&permission, "id = ?", id).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.DatabaseError)
		return nil, appErr
	}

	return permission.ToDto(), nil
}

func (s *repository) Delete(ctx context.Context, payload *dto.PermissionDto) error {
	entity := model.NewPermissionEntity(payload)

	if err := s.db.WithContext(ctx).Delete(entity).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, deleteError), customErrors.DatabaseError)
		return appErr
	}

	return nil
}
