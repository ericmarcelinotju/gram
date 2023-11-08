package setting

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/module/setting"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(group *gin.RouterGroup) func(service setting.Service) {
	settingRoutesFactory := func(service setting.Service) {
		group.GET("", GetSetting(service))
		group.POST("", SaveSetting(service))
	}
	return settingRoutesFactory
}
