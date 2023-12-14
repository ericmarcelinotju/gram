package dto

import (
	"mime/multipart"
	"time"
)

// UserDto struct defines dto of user entity
type UserDto struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string
	Firstname string  `json:"first_name"`
	Lastname  string  `json:"last_name"`
	Title     string  `json:"title"`
	Avatar    *string `json:"avatar,omitempty"`
	RoleId    string  `json:"role_id"`
	RoleName  string  `json:"role_name"`
	Role      RoleDto `json:"role"`

	LastLogin *time.Time `json:"last_login"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostUserDto struct {
	Name            string                `json:"name" form:"name" binding:"required,min=2"`
	Firstname       string                `json:"first_name" form:"first_name"`
	Lastname        string                `json:"last_name" form:"last_name"`
	Title           string                `json:"title" form:"title"`
	Email           string                `json:"email" form:"email" binding:"required,email"`
	Password        string                `json:"password" form:"password" binding:"required"`
	ConfirmPassword string                `json:"confirm_password" form:"confirm_password" binding:"eqfield=Password"`
	Avatar          *multipart.FileHeader `form:"avatar" swaggerignore:"true"`
	RoleId          string                `json:"role_id" form:"role_id" binding:"required,uuid"`
}

type PutUserDto struct {
	Id        string                `json:"-"`
	Name      string                `json:"name" form:"name" binding:"min=2"`
	Firstname string                `json:"first_name" form:"first_name"`
	Lastname  string                `json:"last_name" form:"last_name"`
	Title     string                `json:"title" form:"title"`
	Email     string                `json:"email" form:"email" binding:"email"`
	Avatar    *multipart.FileHeader `form:"avatar" swaggerignore:"true"`
	RoleId    string                `json:"role_id" form:"role_id" binding:"required,uuid"`
}

type GetUserDto struct {
	Name   *string `json:"name" form:"name" uri:"name"`
	Email  *string `json:"email" form:"email" uri:"email"`
	RoleId *string `json:"role_id" form:"role_id" uri:"role_id" binding:"omitempty,uuid"`

	*PaginationDto
	*SortDto
}

type UserChannelDto struct {
	Channel string `json:"channel" form:"channel" uri:"channel"`
}
