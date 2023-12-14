package role

import (
	"net/http"

	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/utils/request"
	"github.com/ericmarcelinotju/gram/utils/response"
	"github.com/gin-gonic/gin"
)

// GetRole godoc
// @Summary     Get list of roles
// @Description Get list of roles
// @Tags        Role
// @Accept      json
// @Produce     json
// @Param       item   query      dto.GetRoleDto   true   "Paging, Search & Filter"
// @Success     200    {object}   response.SetResponse{data=dto.ListRoleDto}
// @Router      /role  [get]
// @Security    Auth
func Get(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.GetRoleDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		roles, total, err := service.Read(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListDto[dto.RoleDto]{
			Data:  roles,
			Total: total,
		}

		response.ResponseSuccess(c, result)
	}
}

// GetRoleDetail godoc
// @Summary     Get role's detail
// @Description  Get role's detail
// @Tags        Role
// @Accept      json
// @Produce     json
// @Param       id          path       string   true   "Role ID"
// @Success     200         {object}   response.SetResponse{data=dto.RoleDto}
// @Router      /role/{id}  [get]
// @Security    Auth
func GetDetail(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := request.BindId(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		result, err := service.ReadById(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, result)
	}
}

// PostRole godoc
// @Summary     Post new role
// @Description Create new role
// @Tags        Role
// @Accept      json
// @Produce     json
// @Param       role   body       dto.PostRoleDto   true   "Role Data"
// @Success     200    {object}   response.SetResponse{data=dto.RoleDto}
// @Router      /role  [post]
// @Security    Auth
func Post(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.PostRoleDto](c)
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

// PutRole godoc
// @Summary     Put role
// @Description Update role datas
// @Tags        Role
// @Accept      json
// @Produce     json
// @Param       id     path       string           true   "Role ID"
// @Param       role   body       dto.PutRoleDto   true   "Role Data"
// @Success     200    {object}   response.SetResponse{data=dto.RoleDto}
// @Router      /role/{id} [put]
// @Security    Auth
func Put(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.PutRoleDto](c)
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

// DeleteRole godoc
// @Summary     Delete role by id
// @Description Delete role by id
// @Tags        Role
// @Accept      json
// @Produce     json
// @Param       id    path       string   true   "Role ID"
// @Success     200   {object}   response.SetResponse
// @Router      /role/{id} [delete]
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
