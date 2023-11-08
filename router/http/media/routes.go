package media

import (
	"github.com/ericmarcelinotju/gram/domain/media"
	"github.com/gin-gonic/gin"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func NewRoutesFactory(group *gin.RouterGroup) func(service media.Service) {
	meterRoutesFactory := func(service media.Service) {
		group.POST("", UploadFile(service))
	}
	return meterRoutesFactory
}
