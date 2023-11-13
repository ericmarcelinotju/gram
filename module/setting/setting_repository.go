package setting

import (
	"context"
	"errors"

	pkgErr "github.com/pkg/errors"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/dto"
	customErrors "github.com/ericmarcelinotju/gram/errors"
	"github.com/ericmarcelinotju/gram/model"
	"github.com/ericmarcelinotju/gram/plugins/cache"

	"gorm.io/gorm"
)

const (
	insertError = "Error in inserting new setting"
	updateError = "Error in updating setting"
	deleteError = "Error in deleting setting"
	selectError = "Error in selecting settings in the database"
)

// Repository provides an abstraction on top of the log data source
type Repository interface {
	Save(ctx context.Context, name, value string) error
	Select(ctx context.Context) ([]dto.SettingDto, error)
	SelectByName(ctx context.Context, name string) (string, error)
	Delete(ctx context.Context, name string) error
}

type repository struct {
	db    *gorm.DB
	cache cache.Cache
}

// New creates a new Store struct
func NewRepository(db *gorm.DB, cache cache.Cache) *repository {
	return &repository{db: db, cache: cache}
}

func (s *repository) Save(ctx context.Context, name, value string) error {
	setting := model.SettingEntity{
		Name:  name,
		Value: value,
	}
	query := s.db.WithContext(ctx).Model(&setting).Where("name = ?", name).Updates(&setting)
	if err := query.Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, insertError), customErrors.CacheError)
		return appErr
	}
	if query.RowsAffected == 0 {
		if err := s.db.WithContext(ctx).Create(&setting).Error; err != nil {
			appErr := customErrors.NewAppError(pkgErr.Wrap(err, insertError), customErrors.DatabaseError)
			return appErr
		}
	}

	if err := s.cache.Del(ctx, "setting"); err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, insertError), customErrors.CacheError)
		return appErr
	}

	if err := s.cache.Set(ctx, "setting-"+name, value, config.Get().Cache.DefaultExpiry); err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, insertError), customErrors.CacheError)
		return appErr
	}

	return nil
}

func (s *repository) Select(ctx context.Context) ([]dto.SettingDto, error) {
	var entities []model.SettingEntity

	// Get from cache
	key := "setting"
	err := s.cache.Get(ctx, key, &entities)
	if err == nil && len(entities) != 0 {
		var results = make([]dto.SettingDto, len(entities))
		for i, element := range entities {
			results[i] = *element.ToDto()
		}
		return results, nil
	}

	query := s.db.WithContext(ctx).
		Model(&model.SettingEntity{}).
		Find(&entities)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := customErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), customErrors.NotFoundError)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.DatabaseError)
		return nil, appErr
	}

	var results = make([]dto.SettingDto, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDto()
	}

	if err := s.cache.Set(ctx, key, entities, config.Get().Cache.DefaultExpiry); err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.CacheError)
		return nil, appErr
	}

	return results, nil
}

func (s *repository) SelectByName(ctx context.Context, name string) (string, error) {
	var value string
	// Get from cache
	key := "setting-" + name
	err := s.cache.Get(ctx, key, &value)
	if err == nil {
		return value, nil
	}

	var entity model.SettingEntity
	query := s.db.WithContext(ctx).
		Model(&model.SettingEntity{}).
		Where("name = ?", name).
		First(&entity)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := customErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), customErrors.NotFoundError)
		return "", appErr
	}

	if err := query.Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.DatabaseError)
		return "", appErr
	}

	if err := s.cache.Set(ctx, key, entity.Value, config.Get().Cache.DefaultExpiry); err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.CacheError)
		return "", appErr
	}

	return entity.Value, nil
}

func (s *repository) Delete(ctx context.Context, name string) error {
	if err := s.db.WithContext(ctx).Where("name = ?", name).Delete(&model.SettingEntity{}).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, deleteError), customErrors.DatabaseError)
		return appErr
	}

	if err := s.cache.Del(ctx, "setting"); err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, insertError), customErrors.CacheError)
		return appErr
	}

	if err := s.cache.Del(ctx, "setting-"+name); err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, deleteError), customErrors.CacheError)
		return appErr
	}

	return nil
}
