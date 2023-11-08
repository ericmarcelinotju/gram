package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	pkgErr "github.com/pkg/errors"

	"github.com/ericmarcelinotju/gram/data/notifier"
	"github.com/ericmarcelinotju/gram/data/schema"
	"github.com/ericmarcelinotju/gram/data/storage"
	domainErrors "github.com/ericmarcelinotju/gram/domain/errors"
	"github.com/ericmarcelinotju/gram/domain/model"
	"gorm.io/gorm"
)

const (
	insertError = "Error in inserting new user"
	updateError = "Error in updating user"
	deleteError = "Error in deleting user"
	selectError = "Error in selecting users in the database"
)

type Store struct {
	db           *gorm.DB
	storage      storage.Storage
	pushNotifier notifier.Notifier
}

// New creates a new Store struct
func New(
	db *gorm.DB,
	storage storage.Storage,
) *Store {
	return &Store{
		db:      db,
		storage: storage,
	}
}

func (s *Store) Migrate(ctx context.Context) error {
	return s.db.AutoMigrate(&schema.User{})
}

func (s *Store) InsertUser(ctx context.Context, payload *model.User) error {
	entity := schema.NewUserSchema(payload)

	if err := s.db.WithContext(ctx).Create(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, insertError), domainErrors.DatabaseError)
		return appErr
	}
	payload.ID = entity.ID.String()
	return nil
}

func (s *Store) UpdateUser(ctx context.Context, payload *model.User) error {
	entity := schema.NewUserSchema(payload)

	if err := s.db.WithContext(ctx).Model(entity).Updates(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, updateError), domainErrors.DatabaseError)
		return appErr
	}
	payload.ID = entity.ID.String()
	return nil
}

func (s *Store) SelectUser(ctx context.Context, filter *model.User) ([]model.User, int64, error) {
	var total int64
	var entities []schema.User

	query := s.db.WithContext(ctx).
		Model(&schema.User{}).
		Preload("Role").
		Preload("Role.Permissions")

	if filter != nil {
		query.Where(schema.NewUserSchema(filter))
	}
	query.Count(&total)
	if filter != nil {
		filter.Pagination.Apply(query)
		filter.Sort.Apply(query)
	}
	query.Find(&entities)

	if err := query.Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, selectError), domainErrors.DatabaseError)
		return nil, total, appErr
	}

	var results = make([]model.User, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDomainModel()
	}

	return results, total, nil
}

func (s *Store) SelectUserByID(ctx context.Context, id string) (*model.User, error) {
	var result schema.User
	query := s.db.
		WithContext(ctx).
		Preload("Role").
		Preload("Role.Permissions").
		First(&result, "id = ?", id)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), domainErrors.NotFoundError)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		return nil, err
	}
	return result.ToDomainModel(), nil
}

func (s *Store) DeleteUser(ctx context.Context, payload *model.User) error {
	entity := schema.NewUserSchema(payload)

	if err := s.db.WithContext(ctx).First(entity).Delete(entity).Error; err != nil {
		appErr := domainErrors.NewAppError(pkgErr.Wrap(err, deleteError), domainErrors.DatabaseError)
		return appErr
	}
	payload.Avatar = entity.Avatar

	return nil
}

func (s *Store) SaveAvatar(payload *model.User) error {
	if payload.AvatarFile == nil {
		appErr := domainErrors.NewAppError(errors.New("uploaded file empty"), domainErrors.ValidationError)
		return appErr
	}
	if payload.Avatar != nil {
		s.storage.Remove(*payload.Avatar)
	}
	filename := "user/" + fmt.Sprintf("%d-%s", time.Now().Unix(), payload.ID)
	if err := s.storage.Upload(*payload.AvatarFile, filename); err != nil {
		appErr := domainErrors.NewAppError(fmt.Errorf("upload file failed: %w", err), domainErrors.RepositoryError)
		return appErr
	}
	payload.Avatar = &filename
	return nil
}

func (s *Store) RemoveAvatar(payload *model.User) error {
	if payload.Avatar != nil {
		return s.storage.Remove(*payload.Avatar)
	}
	return nil
}

func (s *Store) SubscribePushLog(payload *model.User) error {
	if s.pushNotifier == nil {
		return domainErrors.NewAppError(errors.New("no push notification registered"), domainErrors.RepositoryError)
	}
	//TODO :: implement subscribe push notification function
	return domainErrors.NewAppError(errors.New("no notification token supplied"), domainErrors.ValidationError)
}

func (s *Store) UnsubscribePushLog(payload *model.User) error {
	if s.pushNotifier == nil {
		return domainErrors.NewAppError(errors.New("no push notification registered"), domainErrors.RepositoryError)
	}
	//TODO :: implement subscribe push notification function
	return domainErrors.NewAppError(errors.New("no notification token supplied"), domainErrors.ValidationError)
}
