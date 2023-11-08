package permission

import (
	"net/http"

	"github.com/ericmarcelinotju/gram/domain/module/permission"
	dto "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	response "github.com/ericmarcelinotju/gram/utils/http"
	"github.com/gin-gonic/gin"
)

// GetPermission godoc
// @Summary     Get list of permissions
// @Description Get list of permissions
// @Tags        Permission
// @Accept      json
// @Produce     json
// @Param       item   query      PermissionFilter   true   "Paging, Search & Filter"
// @Success     200    {object}   response.SetResponse{data=dto.ListPermissionResponse}
// @Router      /permission  [get]
// @Security    Auth
func GetPermission(service permission.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := BindGet(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		permissions, total, err := service.ReadPermission(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListPermissionResponse{
			Permissions: make([]dto.PermissionResponse, len(permissions)),
			Total:       total,
		}

		for i, permission := range permissions {
			result.Permissions[i] = *permission.ToResponseModel()
		}

		response.ResponseSuccess(c, result)
	}
}

// GetPermissionDetail godoc
// @Summary     Get permission's detail
// @Description  Get permission's detail
// @Tags        Permission
// @Accept      json
// @Produce     json
// @Param       id    path       string   true   "Permission ID"
// @Success     200   {object}   response.SetResponse{data=dto.PermissionResponse}
// @Router      /permission/{id}  [get]
// @Security    Auth
func GetPermissionDetail(service permission.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := BindID(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		permission, err := service.ReadPermissionByID(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, permission.ToResponseModel())
	}
}

// PostPermission godoc
// @Summary     Post new permission
// @Description Create new permission
// @Tags        Permission
// @Accept      json
// @Produce     json
// @Param       permission   body       PostPermissionPayload   true   "Permission Data"
// @Success     200          {object}   response.SetResponse{data=dto.PermissionResponse}
// @Router      /permission  [post]
// @Security    Auth
func PostPermission(service permission.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		permission, err := BindPost(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.CreatePermission(c, permission)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, permission.ToResponseModel())
	}
}

// PutPermission godoc
// @Summary     Put permission
// @Description Update permission datas
// @Tags        Permission
// @Accept      json
// @Produce     json
// @Param       id           path       string           true   "Permission ID"
// @Param       permission   body       PutPermissionPayload   true   "Permission Data"
// @Success     200          {object}   response.SetResponse{data=dto.PermissionResponse}
// @Router      /permission/{id} [put]
// @Security    Auth
func PutPermission(service permission.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		permission, err := BindPut(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.UpdatePermission(c, permission)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, permission.ToResponseModel())
	}
}

// DeletePermission godoc
// @Summary     Delete permission by id
// @Description Delete permission by id
// @Tags        Permission
// @Accept      json
// @Produce     json
// @Param       id    path       string   true   "Permission ID"
// @Success     200   {object}   response.SetResponse
// @Router      /permission/{id} [delete]
// @Security    Auth
func DeletePermission(service permission.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := BindID(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		err = service.DeletePermissionByID(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}
