package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/firelogik/helios/config"
	docs "gitlab.com/firelogik/helios/docs"
)

// NewRoutesFactory create and returns a factory to create routes for the panelment
func Init(group *gin.RouterGroup) func() {

	docs.SwaggerInfo.Title = "GRAM API : Golang Boilerplate"
	docs.SwaggerInfo.Description = "Boilerplate for Golang."
	docs.SwaggerInfo.Version = config.Get().Version
	docs.SwaggerInfo.Host = config.Get().Host.Host
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{config.Get().Host.Scheme}

	routes := func() {
		group.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	return routes
}
