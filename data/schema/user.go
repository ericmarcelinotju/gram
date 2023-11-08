package schema

import (
	"time"

	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/google/uuid"
)

// User struct defines the database model for an user.
type User struct {
	Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string

	Firstname         string
	Lastname          string
	Department        string
	Title             string
	Avatar            *string
	LastLogin         *time.Time
	NotificationToken *string
	RoleID            uuid.UUID
	Role              Role `gorm:"foreignKey:RoleID"`

	ForgotPasswordToken *string
}

func NewUserSchema(entity *model.User) *User {
	id, _ := uuid.Parse(entity.ID)
	roleId, _ := uuid.Parse(entity.RoleID)

	user := &User{
		Model: Model{
			ID:        id,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Username:          entity.Username,
		Email:             entity.Email,
		Password:          entity.Password,
		Firstname:         entity.Firstname,
		Lastname:          entity.Lastname,
		Department:        entity.Department,
		Title:             entity.Title,
		Avatar:            entity.Avatar,
		LastLogin:         entity.LastLogin,
		RoleID:            roleId,
		NotificationToken: entity.NotificationToken,
	}

	return user
}

func (entity *User) ToDomainModel() *model.User {
	user := &model.User{
		ID:         entity.ID.String(),
		Username:   entity.Username,
		Email:      entity.Email,
		Password:   entity.Password,
		Firstname:  entity.Firstname,
		Lastname:   entity.Lastname,
		Department: entity.Department,
		Title:      entity.Title,
		Avatar:     entity.Avatar,
		RoleID:     entity.RoleID.String(),
		Role:       *entity.Role.ToDomainModel(),

		LastLogin: entity.LastLogin,

		NotificationToken: entity.NotificationToken,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return user
}
