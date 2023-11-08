package health

import (
	"time"

	"github.com/gin-gonic/gin"
	response "gitlab.com/firelogik/helios/utils/http"
)

// NewRoutesFactory create and returns a factory to create routes to check API health
func NewRoutesFactory(group *gin.RouterGroup) func() {
	healthRoutesFactory := func() {
		group.GET("", func(c *gin.Context) {
			zone, _ := time.Now().Zone()
			response.ResponseSuccess(c, gin.H{
				"local": time.Local.String(),
				"zone":  zone,
			})
		})
	}

	return healthRoutesFactory
}
