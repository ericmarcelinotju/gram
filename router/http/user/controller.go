package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/module/user"
	dto "gitlab.com/firelogik/helios/router/http/dto/responses"
	response "gitlab.com/firelogik/helios/utils/http"
)

// GetUser godoc
// @Summary     Get list of users
// @Description Get list of users
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       item   query      UserFilter   true   "Paging, Search & Filter"
// @Success     200    {object}   response.SetResponse{data=dto.ListUserResponse}
// @Router      /user  [get]
// @Security    Auth
func GetUser(service user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := BindGet(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		users, total, err := service.ReadUser(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListUserResponse{
			Users: make([]dto.UserResponse, len(users)),
			Total: total,
		}
		for i, user := range users {
			result.Users[i] = *user.ToResponseModel()
		}

		response.ResponseSuccess(c, result)
	}
}

// GetUserDetail godoc
// @Summary     Get user's detail
// @Description  Get user's detail
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id          path       string   true   "User ID"
// @Success     200         {object}   response.SetResponse{data=dto.UserResponse}
// @Router      /user/{id}  [get]
// @Security    Auth
func GetUserDetail(service user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := BindID(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		user, err := service.ReadUserByID(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, user.ToResponseModel())
	}
}

// PostUser godoc
// @Summary     Post new user
// @Description Create new user
// @Tags        User
// @Accept      mpfd
// @Produce     json
// @Param       user   body       PostUserPayload   true   "User Data"
// @Success     200    {object}   response.SetResponse{data=dto.UserResponse}
// @Router      /user  [post]
// @Security    Auth
func PostUser(service user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := BindPost(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.CreateUser(c, user)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, user.ToResponseModel())
	}
}

// PutUser godoc
// @Summary     Put user
// @Description Update user datas
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id		path		string			true	"User ID"
// @Param       user	body		PutUserPayload	true	"User Data"
// @Success     200		{object}	response.SetResponse{data=dto.UserResponse}
// @Router      /user/{id} [put]
// @Security    Auth
func PutUser(service user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := BindPut(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.UpdateUser(c, user)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, user.ToResponseModel())
	}
}

// DeleteUser godoc
// @Summary     Delete user by id
// @Description Delete user by id
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id    path       string   true   "User ID"
// @Success     200   {object}   response.SetResponse
// @Router      /user/{id} [delete]
// @Security    Auth
func DeleteUser(service user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := BindID(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		err = service.DeleteUserByID(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}
