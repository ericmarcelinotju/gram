package role

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/module/role"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(group *gin.RouterGroup) func(service role.Service) {
	roleRoutesFactory := func(service role.Service) {
		group.GET("", GetRole(service))
		group.GET("/:id", GetRoleDetail(service))
		group.POST("", PostRole(service))
		group.PUT("/:id", PutRole(service))
		group.DELETE("/:id", DeleteRole(service))
	}
	return roleRoutesFactory
}
