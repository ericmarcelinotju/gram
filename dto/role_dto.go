package dto

import "time"

// ListRoleDto struct defines http response of roles
type ListRoleDto struct {
	Roles []RoleDto `json:"roles"`
	Total int64     `json:"total"`
}

// RoleDto struct defines dto for role entity
type RoleDto struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permissions []PermissionDto `json:"permissions"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *time.Time      `json:"deleted_at"`
}

type PostRoleDto struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Permissions []IdDto `json:"permissions"`
}

type PutRoleDto struct {
	Id          string  `json:"id" form:"id" uri:"id" binding:"required,uuid"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Permissions []IdDto `json:"permissions"`
}

type GetRoleDto struct {
	Name *string `json:"name" form:"name"`

	*PaginationDto
	*SortDto
}
