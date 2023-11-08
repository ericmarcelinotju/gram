package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/module/role"
	dto "gitlab.com/firelogik/helios/router/http/dto/responses"
	response "gitlab.com/firelogik/helios/utils/http"
)

// GetRole godoc
// @Summary     Get list of roles
// @Description Get list of roles
// @Tags        Role
// @Accept      json
// @Produce     json
// @Param       item   query      RoleFilter   true   "Paging, Search & Filter"
// @Success     200    {object}   response.SetResponse{data=dto.ListRoleResponse}
// @Router      /role  [get]
// @Security    Auth
func GetRole(service role.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := BindGet(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		roles, total, err := service.ReadRole(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListRoleResponse{
			Roles: make([]dto.RoleResponse, len(roles)),
			Total: total,
		}

		for i, role := range roles {
			result.Roles[i] = *role.ToResponseModel()
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
// @Success     200         {object}   response.SetResponse{data=dto.RoleResponse}
// @Router      /role/{id}  [get]
// @Security    Auth
func GetRoleDetail(service role.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := BindID(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		result, err := service.ReadRoleByID(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, result.ToResponseModel())
	}
}

// PostRole godoc
// @Summary     Post new role
// @Description Create new role
// @Tags        Role
// @Accept      json
// @Produce     json
// @Param       role   body       PostRolePayload   true   "Role Data"
// @Success     200    {object}   response.SetResponse{data=dto.RoleResponse}
// @Router      /role  [post]
// @Security    Auth
func PostRole(service role.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		role, err := BindPost(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.CreateRole(c, role)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, role.ToResponseModel())
	}
}

// PutRole godoc
// @Summary     Put role
// @Description Update role datas
// @Tags        Role
// @Accept      json
// @Produce     json
// @Param       id     path       string           true   "Role ID"
// @Param       role   body       PutRolePayload   true   "Role Data"
// @Success     200    {object}   response.SetResponse{data=dto.RoleResponse}
// @Router      /role/{id} [put]
// @Security    Auth
func PutRole(service role.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		role, err := BindPut(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.UpdateRole(c, role)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, role.ToResponseModel())
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
func DeleteRole(service role.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := BindID(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		err = service.DeleteRoleByID(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}
