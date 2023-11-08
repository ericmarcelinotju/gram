package permission

import (
	"github.com/ericmarcelinotju/gram/domain/module/permission"
	"github.com/gin-gonic/gin"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(group *gin.RouterGroup) func(service permission.Service) {
	permissionRoutesFactory := func(service permission.Service) {
		group.GET("", GetPermission(service))
		group.GET("/:id", GetPermissionDetail(service))
		group.POST("", PostPermission(service))
		group.PUT("/:id", PutPermission(service))
		group.DELETE("/:id", DeletePermission(service))
	}
	return permissionRoutesFactory
}
