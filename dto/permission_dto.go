package dto

import (
	"time"
)

// ListPermissionDto struct defines http response of permissions
type ListPermissionDto struct {
	Permissions []PermissionDto `json:"permissions"`
	Total       int64           `json:"total"`
}

// PermissionDto struct defines dto for permission entity
type PermissionDto struct {
	Id          string     `json:"id"`
	Method      string     `json:"method"`
	Module      string     `json:"module"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type PostPermissionDto struct {
	Method      string `json:"method" binding:"required"`
	Module      string `json:"module" binding:"required"`
	Description string `json:"description"`
}

type PutPermissionDto struct {
	Id          string `json:"id" form:"id" uri:"id" binding:"required,uuid"`
	Method      string `json:"method"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

type GetPermissionDto struct {
	Method *string `json:"method" form:"method"`
	Module *string `json:"module" form:"module"`

	*PaginationDto
	*SortDto
}
