package permission

import (
	"github.com/gin-gonic/gin"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(router *gin.RouterGroup) func(service Service) {
	group := router.Group("/api/permission")
	permissionRoutesFactory := func(service Service) {
		group.GET("", Get(service))
		group.GET("/:id", GetDetail(service))
		group.POST("", Post(service))
		group.PUT("/:id", Put(service))
		group.DELETE("/:id", Delete(service))
	}
	return permissionRoutesFactory
}
