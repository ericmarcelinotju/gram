package permission

import (
	"net/http"

	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/utils/request"
	"github.com/ericmarcelinotju/gram/utils/response"
	"github.com/gin-gonic/gin"
)

// GetPermission godoc
// @Summary     Get list of permissions
// @Description Get list of permissions
// @Tags        Permission
// @Accept      json
// @Produce     json
// @Param       item   query      dto.GetPermissionDto   true   "Paging, Search & Filter"
// @Success     200    {object}   response.SetResponse{data=dto.ListPermissionDto}
// @Router      /permission  [get]
// @Security    Auth
func Get(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.GetPermissionDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		permissions, total, err := service.Read(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListDto[dto.PermissionDto]{
			Data:  permissions,
			Total: total,
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
// @Success     200   {object}   response.SetResponse{data=dto.PermissionDto}
// @Router      /permission/{id}  [get]
// @Security    Auth
func GetDetail(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := request.BindId(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		res, err := service.ReadById(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, res)
	}
}

// PostPermission godoc
// @Summary     Post new permission
// @Description Create new permission
// @Tags        Permission
// @Accept      json
// @Produce     json
// @Param       permission   body       dto.PostPermissionDto   true   "Permission Data"
// @Success     200          {object}   response.SetResponse{data=dto.PermissionDto}
// @Router      /permission  [post]
// @Security    Auth
func Post(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.PostPermissionDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		res, err := service.Create(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, res)
	}
}

// PutPermission godoc
// @Summary     Put permission
// @Description Update permission datas
// @Tags        Permission
// @Accept      json
// @Produce     json
// @Param       id           path       string           			 true   "Permission ID"
// @Param       permission   body       dto.PutPermissionDto   true   "Permission Data"
// @Success     200          {object}   response.SetResponse{data=dto.PermissionDto}
// @Router      /permission/{id} [put]
// @Security    Auth
func Put(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.PutPermissionDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		id, err := request.BindId(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		payload.Id = id

		res, err := service.Update(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, res)
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
func Delete(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := request.BindId(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		err = service.DeleteById(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}
