package permission

import (
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/ericmarcelinotju/gram/router/http/dto"
	"github.com/gin-gonic/gin"
)

type PostPermissionPayload struct {
	Method      string `json:"method" binding:"required"`
	Module      string `json:"module" binding:"required"`
	Description string `json:"description"`
}

type PutPermissionPayload struct {
	Method      string `json:"method"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

type PermissionFilter struct {
	ID     *string `json:"id" form:"id" uri:"id" binding:"omitempty,uuid"`
	Method *string `json:"method" form:"method"`
	Module *string `json:"module" form:"module"`

	dto.Pagination
	dto.Sort
}

func BindPost(c *gin.Context) (*model.Permission, error) {
	var json PostPermissionPayload
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	permission := &model.Permission{
		Method:      json.Method,
		Module:      json.Module,
		Description: json.Description,
	}

	return permission, nil
}

func BindPut(c *gin.Context) (*model.Permission, error) {
	var idParam dto.IdParam
	err := c.ShouldBindUri(&idParam)
	if err != nil {
		return nil, err
	}
	var json PutPermissionPayload
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	payload := &model.Permission{
		ID:          idParam.Id,
		Method:      json.Method,
		Module:      json.Module,
		Description: json.Description,
	}

	return payload, nil
}

func BindGet(c *gin.Context) (*model.Permission, error) {
	var json PermissionFilter
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	var permission *model.Permission = &model.Permission{}

	if json.ID != nil {
		permission.ID = *json.ID
	}
	if json.Method != nil {
		permission.Method = *json.Method
	}
	if json.Module != nil {
		permission.Module = *json.Module
	}
	pagination, _ := dto.ToPaginationModel(json.Pagination)
	if pagination != nil {
		permission.Pagination = *pagination
	}
	sort, _ := dto.ToSortModel(json.Sort)
	if sort != nil {
		permission.Sort = *sort
	}

	return permission, nil
}

func BindID(c *gin.Context) (string, error) {
	var json PermissionFilter
	if err := c.ShouldBindUri(&json); err != nil {
		return "", err
	}

	return *json.ID, nil
}
