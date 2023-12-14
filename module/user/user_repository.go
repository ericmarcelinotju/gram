package user

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"gorm.io/gorm"

	"github.com/gorilla/websocket"
	pkgErr "github.com/pkg/errors"

	"github.com/ericmarcelinotju/gram/dto"
	customErrors "github.com/ericmarcelinotju/gram/errors"
	"github.com/ericmarcelinotju/gram/model"
	"github.com/ericmarcelinotju/gram/plugins/storage"
	ws "github.com/ericmarcelinotju/gram/plugins/websocket"
	"github.com/ericmarcelinotju/gram/utils/crypt"
)

const (
	insertError = "Error in inserting new user"
	updateError = "Error in updating user"
	deleteError = "Error in deleting user"
	selectError = "Error in selecting users in the database"
)

// Repository provides an abstraction on top of the user data source
type Repository interface {
	Insert(context.Context, *dto.UserDto) error
	Update(context.Context, *dto.UserDto) error
	UpdatePassword(ctx context.Context, id string, password string) error
	Select(context.Context, *dto.UserDto, *dto.PaginationDto, *dto.SortDto) ([]dto.UserDto, int64, error)
	SelectById(context.Context, string) (*dto.UserDto, error)
	SelectByUsername(context.Context, string) (*dto.UserDto, error)
	Delete(context.Context, *dto.UserDto) error

	SaveAvatar(file *multipart.File, filename string) error
	RemoveAvatar(filename string) error

	Connect(*websocket.Conn, *dto.UserChannelDto) error
}

type repository struct {
	db         *gorm.DB
	storage    storage.Storage
	dispatcher *ws.Dispatcher
}

// New creates a new repository struct
func NewRepository(
	db *gorm.DB,
	storage storage.Storage,
	dispatcher *ws.Dispatcher,
) *repository {
	return &repository{
		db:         db,
		storage:    storage,
		dispatcher: dispatcher,
	}
}

func (s *repository) Insert(ctx context.Context, payload *dto.UserDto) error {
	entity := model.NewUserEntity(payload)

	hashedPassword, err := crypt.Hash(entity.Password)
	if err != nil {
		return err
	}
	entity.Password = hashedPassword

	if err := s.db.WithContext(ctx).Create(entity).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, insertError), customErrors.DatabaseError)
		return appErr
	}
	payload.Id = entity.Id.String()
	return nil
}

func (s *repository) Update(ctx context.Context, payload *dto.UserDto) error {
	entity := model.NewUserEntity(payload)

	if err := s.db.WithContext(ctx).Model(entity).Updates(entity).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, updateError), customErrors.DatabaseError)
		return appErr
	}
	payload.Id = entity.Id.String()
	return nil
}

func (s *repository) UpdatePassword(ctx context.Context, id string, password string) error {
	hashedPassword, err := crypt.Hash(password)
	if err != nil {
		return err
	}

	if err := s.db.WithContext(ctx).Model(&model.UserEntity{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, updateError), customErrors.DatabaseError)
		return appErr
	}
	return nil
}

func (s *repository) Select(
	ctx context.Context,
	filter *dto.UserDto,
	pagination *dto.PaginationDto,
	sort *dto.SortDto,
) ([]dto.UserDto, int64, error) {
	var total int64
	var entities []model.UserEntity

	query := s.db.WithContext(ctx).
		Model(&model.UserEntity{}).
		Preload("Role").
		Preload("Role.Permissions")

	if filter != nil {
		query.Where(model.NewUserEntity(filter))
	}
	query.Count(&total)
	if pagination != nil {
		pagination.Apply(query)
	}
	if sort != nil {
		sort.Apply(query)
	}
	query.Find(&entities)

	if err := query.Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.DatabaseError)
		return nil, total, appErr
	}

	var results = make([]dto.UserDto, len(entities))
	for i, element := range entities {
		results[i] = *element.ToDto()
	}

	return results, total, nil
}

func (s *repository) SelectById(ctx context.Context, id string) (*dto.UserDto, error) {
	var result model.UserEntity
	query := s.db.
		WithContext(ctx).
		Preload("Role").
		Preload("Role.Permissions").
		First(&result, "id = ?", id)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := customErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), customErrors.NotFoundError)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		return nil, err
	}
	return result.ToDto(), nil
}

func (s *repository) SelectByUsername(ctx context.Context, id string) (*dto.UserDto, error) {
	var result model.UserEntity
	query := s.db.
		WithContext(ctx).
		Preload("Role").
		Preload("Role.Permissions").
		First(&result, "name = ?", id)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		appErr := customErrors.NewAppError(pkgErr.Wrap(query.Error, selectError), customErrors.NotFoundError)
		return nil, appErr
	}

	if err := query.Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, selectError), customErrors.RepositoryError)
		return nil, appErr
	}
	return result.ToDto(), nil
}

func (s *repository) Delete(ctx context.Context, payload *dto.UserDto) error {
	entity := model.NewUserEntity(payload)

	if err := s.db.WithContext(ctx).First(entity).Delete(entity).Error; err != nil {
		appErr := customErrors.NewAppError(pkgErr.Wrap(err, deleteError), customErrors.DatabaseError)
		return appErr
	}
	payload.Avatar = entity.Avatar

	return nil
}

func (s *repository) SaveAvatar(file *multipart.File, filename string) error {
	if file == nil {
		appErr := customErrors.NewAppError(errors.New("uploaded file empty"), customErrors.ValidationError)
		return appErr
	}
	if err := s.storage.Upload(*file, filename); err != nil {
		appErr := customErrors.NewAppError(fmt.Errorf("upload file failed: %w", err), customErrors.RepositoryError)
		return appErr
	}
	return nil
}

func (s *repository) RemoveAvatar(filename string) error {
	return s.storage.Remove(filename)
}

func (s *repository) Connect(conn *websocket.Conn, channel *dto.UserChannelDto) error {
	client := ws.NewClient(s.dispatcher, conn, channel.Channel)
	s.dispatcher.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WriteDispatch()
	go client.ReadDispatch()

	return nil
}
