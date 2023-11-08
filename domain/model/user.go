package model

import (
	"mime/multipart"
	"os"
	"time"

	dto "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	timeUtils "github.com/ericmarcelinotju/gram/utils/time"
)

// User struct defines the response model for a user APIs.
type User struct {
	ID         string
	Username   string
	Email      string
	Password   string
	Firstname  string
	Lastname   string
	Department string
	Title      string
	RoleID     string
	Role       Role

	LastLogin  *time.Time
	Avatar     *string
	AvatarFile *multipart.File

	NotificationToken *string

	Pagination Pagination

	Sort Sort

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (entity *User) ToResponseModel() *dto.UserResponse {
	if entity.Avatar != nil && *entity.Avatar != "" {
		avatarURL := os.Getenv("MEDIA_URL") + *entity.Avatar
		entity.Avatar = &avatarURL
	}

	response := &dto.UserResponse{
		ID:         entity.ID,
		Username:   entity.Username,
		Email:      entity.Email,
		Firstname:  entity.Firstname,
		Lastname:   entity.Lastname,
		Department: entity.Department,
		Title:      entity.Title,
		Avatar:     entity.Avatar,
		RoleID:     entity.RoleID,
		RoleName:   entity.Role.Name,
		Role:       *entity.Role.ToResponseModel(),
		LastLogin:  timeUtils.FormatResponseTime(entity.LastLogin),

		CreatedAt: *timeUtils.FormatResponseTime(&entity.CreatedAt),
		UpdatedAt: *timeUtils.FormatResponseTime(&entity.UpdatedAt),
	}
	return response
}
