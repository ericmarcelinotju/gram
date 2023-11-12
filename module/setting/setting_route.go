package setting

import (
	"github.com/gin-gonic/gin"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(router *gin.RouterGroup) func(service Service) {
	group := router.Group("/api/setting")
	settingRoutesFactory := func(service Service) {
		group.GET("", Get(service))
		group.POST("", Save(service))
	}
	return settingRoutesFactory
}
