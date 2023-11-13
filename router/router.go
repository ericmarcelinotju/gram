package router

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	authModule "github.com/ericmarcelinotju/gram/module/auth"
	healthModule "github.com/ericmarcelinotju/gram/module/health"
	permissionModule "github.com/ericmarcelinotju/gram/module/permission"
	roleModule "github.com/ericmarcelinotju/gram/module/role"
	settingModule "github.com/ericmarcelinotju/gram/module/setting"
	userModule "github.com/ericmarcelinotju/gram/module/user"
	"github.com/ericmarcelinotju/gram/plugins/job"

	swaggerRoutes "github.com/ericmarcelinotju/gram/router/swagger"

	"github.com/ericmarcelinotju/gram/router/middleware"

	response "github.com/ericmarcelinotju/gram/utils/http"
)

// NewHTTPHandler returns the HTTP requests handler
func NewHTTPHandler(
	authSvc authModule.Service,

	userSvc userModule.Service,
	roleSvc roleModule.Service,
	permissionSvc permissionModule.Service,

	settingSvc settingModule.Service,

	queueConnection rmq.Connection,

	backupQueue *job.Queue,
) http.Handler {

	gin.DefaultWriter = log.Writer()

	router := gin.Default()

	config := cors.Config{
		AllowOrigins: []string{
			"http://localhost:9099",
			"http://localhost:80",
			"http://localhost",
			"http://10.224.171.167:80",
			"http://10.224.171.167",
			"https://10.224.171.167:80",
			"https://10.224.171.167",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-XSRF-TOKEN", "App-Name", "ResponseType"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	router.NoRoute(func(c *gin.Context) {
		response.ResponseError(c, errors.New("route not found"), http.StatusNotFound)
	})

	router.Static("/media", "./media")

	healthGroup := router.Group("/health")
	healthModule.NewRoutesFactory(healthGroup)()

	authModule.NewRoutesFactory(router)(authSvc, userSvc)

	authMiddleware := middleware.NewAuthMiddleware(authSvc)
	authGroup := router.Group("")
	authGroup.Use(authMiddleware.Authenticate)
	authGroup.Use(authMiddleware.Authorize)
	{
		userModule.NewApiRoutesFactory(authGroup)(userSvc)
		roleModule.NewRoutesFactory(authGroup)(roleSvc)
		permissionModule.NewRoutesFactory(authGroup)(permissionSvc)
		settingModule.NewRoutesFactory(authGroup)(settingSvc)
	}

	swaggerRoutes.Init(router.Group("swagger"))()

	return router
}
