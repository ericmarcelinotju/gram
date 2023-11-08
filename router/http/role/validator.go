package role

import (
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/ericmarcelinotju/gram/router/http/dto"
	"github.com/gin-gonic/gin"
)

type PostRolePayload struct {
	Name        string        `json:"name" binding:"required"`
	Description string        `json:"description"`
	Permissions []dto.IdParam `json:"permissions"`
}

type PutRolePayload struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Permissions []dto.IdParam `json:"permissions"`
}

type RoleFilter struct {
	ID   *string `json:"id" form:"id" uri:"id" binding:"omitempty,uuid"`
	Name *string `json:"name" form:"name"`

	dto.Pagination
	dto.Sort
}

func BindPost(c *gin.Context) (*model.Role, error) {
	var json PostRolePayload
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	var permissions []model.Permission = make([]model.Permission, len(json.Permissions))
	for i, item := range json.Permissions {
		permissions[i] = model.Permission{ID: item.Id}
	}

	role := &model.Role{
		Name:        json.Name,
		Description: json.Description,
		Permissions: permissions,
	}

	return role, nil
}

func BindPut(c *gin.Context) (*model.Role, error) {
	var idParam dto.IdParam
	err := c.ShouldBindUri(&idParam)
	if err != nil {
		return nil, err
	}
	var json PutRolePayload
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	var permissions []model.Permission = make([]model.Permission, len(json.Permissions))
	for i, item := range json.Permissions {
		permissions[i] = model.Permission{ID: item.Id}
	}

	payload := &model.Role{
		ID:          idParam.Id,
		Name:        json.Name,
		Description: json.Description,
		Permissions: permissions,
	}
	return payload, nil
}

func BindGet(c *gin.Context) (*model.Role, error) {
	var json RoleFilter
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	var role *model.Role = &model.Role{}

	if json.ID != nil {
		role.ID = *json.ID
	}
	if json.Name != nil {
		role.Name = *json.Name
	}
	pagination, _ := dto.ToPaginationModel(json.Pagination)
	if pagination != nil {
		role.Pagination = *pagination
	}
	sort, _ := dto.ToSortModel(json.Sort)
	if sort != nil {
		role.Sort = *sort
	}

	return role, nil
}

func BindID(c *gin.Context) (string, error) {
	var json RoleFilter
	if err := c.ShouldBindUri(&json); err != nil {
		return "", err
	}

	return *json.ID, nil
}
