package setting

import (
	"github.com/ericmarcelinotju/gram/domain/module/setting"
	"github.com/gin-gonic/gin"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(group *gin.RouterGroup) func(service setting.Service) {
	settingRoutesFactory := func(service setting.Service) {
		group.GET("", GetSetting(service))
		group.POST("", SaveSetting(service))
	}
	return settingRoutesFactory
}
