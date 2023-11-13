package user

import (
	"github.com/gin-gonic/gin"
)

// NewApiRoutesFactory create and returns a factory to create routes for the panelment
func NewApiRoutesFactory(router *gin.RouterGroup) func(service Service) {
	group := router.Group("/api/user")
	userRoutesFactory := func(service Service) {
		group.GET("", Get(service))
		group.GET("/:id", GetDetail(service))
		group.POST("", Post(service))
		group.PUT("/:id", Put(service))
		group.DELETE("/:id", Delete(service))
	}
	return userRoutesFactory
}

// NewWsRoutesFactory create and returns a factory to create routes for the panelment
func NewWsRoutesFactory(router *gin.RouterGroup) func(service Service) {
	group := router.Group("/ws/user")
	websocketRoutesFactory := func(service Service) {
		group.GET("/:channel", Connect(service))
	}

	return websocketRoutesFactory
}
