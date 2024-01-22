package model

import (
	uuid "github.com/satori/go.uuid"
	"time"

	"github.com/ericmarcelinotju/gram/dto"
)

// UserEntity struct defines the database model for an user.
type UserEntity struct {
	Model
	Name     string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string

	Firstname string
	Lastname  string
	Title     string
	Avatar    *string
	LastLogin *time.Time
	RoleId    uuid.UUID
	Role      RoleEntity `gorm:"foreignKey:RoleId"`

	ForgotPasswordToken *string
}

func (UserEntity) TableName() string {
	return "users"
}

func NewUserEntity(entity *dto.UserDto) *UserEntity {
	id, _ := uuid.FromString(entity.Id)
	roleId, _ := uuid.FromString(entity.RoleId)

	user := &UserEntity{
		Model: Model{
			Id:        id,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Name:      entity.Name,
		Email:     entity.Email,
		Password:  entity.Password,
		Firstname: entity.Firstname,
		Lastname:  entity.Lastname,
		Title:     entity.Title,
		Avatar:    entity.Avatar,
		LastLogin: entity.LastLogin,
		RoleId:    roleId,
	}

	return user
}

func (entity *UserEntity) ToDto() *dto.UserDto {
	user := &dto.UserDto{
		Id:        entity.Id.String(),
		Name:      entity.Name,
		Email:     entity.Email,
		Password:  entity.Password,
		Firstname: entity.Firstname,
		Lastname:  entity.Lastname,
		Title:     entity.Title,
		Avatar:    entity.Avatar,
		RoleId:    entity.RoleId.String(),
		Role:      *entity.Role.ToDto(),

		LastLogin: entity.LastLogin,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return user
}
