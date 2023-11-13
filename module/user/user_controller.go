package user

import (
	"net/http"

	"github.com/ericmarcelinotju/gram/dto"
	httpUtil "github.com/ericmarcelinotju/gram/utils/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
func Get(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := httpUtil.Bind[dto.GetUserDto](c)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		users, total, err := service.Read(c, payload)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListUserDto{
			Users: users,
			Total: total,
		}

		httpUtil.ResponseSuccess(c, result)
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
func GetDetail(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := httpUtil.BindId(c)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		user, err := service.ReadById(c, id)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		httpUtil.ResponseSuccess(c, user)
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
func Post(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := httpUtil.Bind[dto.PostUserDto](c)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		res, err := service.Create(c, payload)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		httpUtil.ResponseSuccess(c, res)
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
func Put(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := httpUtil.Bind[dto.PutUserDto](c)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		res, err := service.Update(c, payload)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		httpUtil.ResponseSuccess(c, res)
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
func Delete(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := httpUtil.BindId(c)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		err = service.DeleteById(c, id)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		httpUtil.ResponseSuccess(c, nil)
	}
}

func Connect(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		channel, err := httpUtil.Bind[dto.UserChannelDto](c)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		// upgrader upgrades the request to WS
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}
		// serveWs handles websocket requests from the peer.
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		err = service.Connect(conn, channel)
		if err != nil {
			return
		}
	}
}
