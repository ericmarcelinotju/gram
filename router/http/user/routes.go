package user

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/module/user"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(group *gin.RouterGroup) func(service user.Service) {
	userRoutesFactory := func(service user.Service) {
		group.GET("", GetUser(service))
		group.GET("/:id", GetUserDetail(service))
		group.POST("", PostUser(service))
		group.PUT("/:id", PutUser(service))
		group.DELETE("/:id", DeleteUser(service))
	}
	return userRoutesFactory
}
