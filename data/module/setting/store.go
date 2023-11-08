package unit

import (
	"context"
	"errors"

	pkgErr "github.com/pkg/errors"

	"gitlab.com/firelogik/helios/config"
	"gitlab.com/firelogik/helios/data/cache"
	"gitlab.com/firelogik/helios/data/schema"
	domainErrors "gitlab.com/firelogik/helios/domain/errors"
	"gitlab.com/firelogik/helios/domain/model"

	"gorm.io/gorm"
)

const (
	insertError = "Error in inserting new setting"
	updateError = "Error in updating setting"
	deleteError = "Error in deleting setting"
	selectError = "Error in selecting settings in the database"
)

type Store struct {
	db    *gorm.DB
	cache cache.Cache
}

// New creates a new Store struct
func New(db *gorm.DB, cache cache.Cache) *Store {
	return &Store{db: db, cache: cache}
}

func (s *Store) Migrate(ctx context.Context) error {
	return s.db.AutoMigrate(&schema.Setting{})
}

func (s *Store) SaveSetting(ctx context.Context, name, value string) error {
	setting := schema.Setting{
		Name:  name,
		Value: value,
	}
	query := s.db.WithContext(ctx).Model(&setting).Where("name = ?", name).Updates(&setting)
	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.CacheError)
		return appErr
	}
	if query.RowsAffected == 0 {
		if err := s.db.WithContext(ctx).Create(&setting).Error; err != nil {
			appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.DatabaseError)
			return appErr
		}
	}

	if err := s.cache.Del(ctx, "setting"); err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.CacheError)
		return appErr
	}

	if err := s.cache.Set(ctx, "setting-"+name, value, config.Get().Cache.DefaultExpiry); err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.CacheError)
		return appErr
	}

	return nil
}

func (s *Store) SelectSetting(ctx context.Context) ([]model.Setting, error) {
	var entities []schema.Setting

	// Get from cache
	key := "setting"
	err := s.cache.Get(ctx, key, &entities)
	if err == nil && len(entities) != 0 {
		var results = make([]model.Setting, len(entities))
		for i, element := range entities {
			results[i] = *element.ToDomainModel()
		}
		return results, nil
	}

	query := s.db.WithContext(ctx).
		Model(&schema.Setting{}).
		Find(&entities)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), domainErrors.NotFoundError)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.DatabaseError)
		return nil, appErr
	}

	var results = make([]model.Setting, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDomainModel()
	}

	if err := s.cache.Set(ctx, key, entities, config.Get().Cache.DefaultExpiry); err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.CacheError)
		return nil, appErr
	}

	return results, nil
}

func (s *Store) SelectSettingByName(ctx context.Context, name string) (string, error) {
	var value string
	// Get from cache
	key := "setting-" + name
	err := s.cache.Get(ctx, key, &value)
	if err == nil {
		return value, nil
	}

	var entity schema.Setting
	query := s.db.WithContext(ctx).
		Model(&schema.Setting{}).
		Where("name = ?", name).
		First(&entity)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), domainErrors.NotFoundError)
		return "", appErr
	}

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.DatabaseError)
		return "", appErr
	}

	if err := s.cache.Set(ctx, key, entity.Value, config.Get().Cache.DefaultExpiry); err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.CacheError)
		return "", appErr
	}

	return entity.Value, nil
}

func (s *Store) DeleteSetting(ctx context.Context, name string) error {
	if err := s.db.WithContext(ctx).Where("name = ?", name).Delete(&schema.Setting{}).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, deleteError), domainErrors.DatabaseError)
		return appErr
	}

	if err := s.cache.Del(ctx, "setting"); err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.CacheError)
		return appErr
	}

	if err := s.cache.Del(ctx, "setting-"+name); err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, deleteError), domainErrors.CacheError)
		return appErr
	}

	return nil
}
