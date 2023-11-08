package log

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/module/log"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(group *gin.RouterGroup) func(service log.Service) {
	logRoutesFactory := func(service log.Service) {
		group.GET("", GetLog(service))
		group.GET("/:id", GetLogDetail(service))
		group.POST("", PostLog(service))
		group.DELETE("/:id", DeleteLog(service))
	}
	return logRoutesFactory
}
