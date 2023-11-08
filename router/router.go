package router

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/ericmarcelinotju/gram/data/job"
	"github.com/ericmarcelinotju/gram/domain/media"
	"github.com/ericmarcelinotju/gram/domain/module/auth"
	logService "github.com/ericmarcelinotju/gram/domain/module/log"
	"github.com/ericmarcelinotju/gram/domain/module/permission"
	"github.com/ericmarcelinotju/gram/domain/module/role"
	"github.com/ericmarcelinotju/gram/domain/module/setting"
	"github.com/ericmarcelinotju/gram/domain/module/user"
	"github.com/ericmarcelinotju/gram/domain/websocket"
	healthRoutes "github.com/ericmarcelinotju/gram/router/http/health"

	authRoutes "github.com/ericmarcelinotju/gram/router/http/auth"
	mediaRoutes "github.com/ericmarcelinotju/gram/router/http/media"
	permissionRoutes "github.com/ericmarcelinotju/gram/router/http/permission"
	roleRoutes "github.com/ericmarcelinotju/gram/router/http/role"
	userRoutes "github.com/ericmarcelinotju/gram/router/http/user"

	logRoutes "github.com/ericmarcelinotju/gram/router/http/log"
	settingRoutes "github.com/ericmarcelinotju/gram/router/http/setting"

	swaggerRoutes "github.com/ericmarcelinotju/gram/router/http/swagger"
	websocketRoutes "github.com/ericmarcelinotju/gram/router/websocket"

	"github.com/ericmarcelinotju/gram/router/middleware"

	response "github.com/ericmarcelinotju/gram/utils/http"
)

// NewHTTPHandler returns the HTTP requests handler
func NewHTTPHandler(
	authSvc auth.Service,
	mediaSvc media.Service,

	userSvc user.Service,
	roleSvc role.Service,
	permissionSvc permission.Service,

	logSvc logService.Service,
	settingSvc setting.Service,

	websocketSvc websocket.WebsocketService,

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

	apiGroup := router.Group("/api")

	healthGroup := apiGroup.Group("/health")
	healthRoutes.NewRoutesFactory(healthGroup)()

	authGroup := apiGroup.Group("/auth")
	authRoutes.NewRoutesFactory(authGroup)(authSvc, userSvc)

	authMiddleware := middleware.NewAuthMiddleware(authSvc)

	apiGroup.Use(authMiddleware.Authenticate)
	{
		apiGroup.Use(authMiddleware.Authorize)
		{
			userGroup := apiGroup.Group("/user")
			userRoutes.NewRoutesFactory(userGroup)(userSvc)

			roleGroup := apiGroup.Group("/role")
			roleRoutes.NewRoutesFactory(roleGroup)(roleSvc)

			permissionGroup := apiGroup.Group("/permission")
			permissionRoutes.NewRoutesFactory(permissionGroup)(permissionSvc)

			logGroup := apiGroup.Group("/log")
			logRoutes.NewRoutesFactory(logGroup)(logSvc)

			settingGroup := apiGroup.Group("/setting")
			settingRoutes.NewRoutesFactory(settingGroup)(settingSvc)
		}
	}

	mediaGroup := router.Group("/media")
	mediaRoutes.NewRoutesFactory(mediaGroup)(mediaSvc)

	wsGroup := router.Group("/ws")
	websocketRoutes.NewRoutesFactory(wsGroup)(websocketSvc)

	swaggerRoutes.Init(router.Group("swagger"))()

	return router
}
