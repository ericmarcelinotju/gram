package dto

import (
	"mime/multipart"
	"time"
)

// ListUserDto struct defines http response of users
type ListUserDto struct {
	Users []UserDto `json:"users"`
	Total int64     `json:"total"`
}

// UserDto struct defines dto of user entity
type UserDto struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string
	Firstname  string  `json:"first_name"`
	Lastname   string  `json:"last_name"`
	Department string  `json:"department"`
	Title      string  `json:"title"`
	Avatar     *string `json:"avatar,omitempty"`
	RoleId     string  `json:"role_id"`
	RoleName   string  `json:"role_name"`
	Role       RoleDto `json:"role"`

	LastLogin *time.Time `json:"last_login"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostUserDto struct {
	Username        string                `json:"username" form:"username" binding:"required,min=2"`
	Firstname       string                `json:"first_name" form:"first_name"`
	Lastname        string                `json:"last_name" form:"last_name"`
	Department      string                `json:"department" form:"department"`
	Title           string                `json:"title" form:"title"`
	Email           string                `json:"email" form:"email" binding:"required,email"`
	Password        string                `json:"password" form:"password" binding:"required"`
	ConfirmPassword string                `json:"confirm_password" form:"confirm_password" binding:"eqfield=Password"`
	Avatar          *multipart.FileHeader `form:"avatar" swaggerignore:"true"`
	RoleId          string                `json:"role_id" form:"role_id" binding:"required,uuid"`
}

type PutUserDto struct {
	Id         string                `json:"id" form:"id" uri:"id" binding:"required,uuid"`
	Username   string                `json:"username" form:"username" binding:"min=2"`
	Firstname  string                `json:"first_name" form:"first_name"`
	Lastname   string                `json:"last_name" form:"last_name"`
	Department string                `json:"department" form:"department"`
	Title      string                `json:"title" form:"title"`
	Email      string                `json:"email" form:"email" binding:"email"`
	Avatar     *multipart.FileHeader `form:"avatar" swaggerignore:"true"`
	RoleId     string                `json:"role_id" form:"role_id" binding:"required,uuid"`
}

type GetUserDto struct {
	Username *string `json:"username" form:"username" uri:"username"`
	Email    *string `json:"email" form:"email" uri:"email"`
	RoleId   *string `json:"role_id" form:"role_id" uri:"role_id" binding:"omitempty,uuid"`

	*PaginationDto
	*SortDto
}

type UserChannelDto struct {
	Channel string `json:"channel" form:"channel" uri:"channel"`
}
