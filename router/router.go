package router

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gitlab.com/firelogik/helios/data/job"
	"gitlab.com/firelogik/helios/domain/media"
	"gitlab.com/firelogik/helios/domain/module/auth"
	logService "gitlab.com/firelogik/helios/domain/module/log"
	"gitlab.com/firelogik/helios/domain/module/permission"
	"gitlab.com/firelogik/helios/domain/module/role"
	"gitlab.com/firelogik/helios/domain/module/setting"
	"gitlab.com/firelogik/helios/domain/module/user"
	"gitlab.com/firelogik/helios/domain/websocket"
	healthRoutes "gitlab.com/firelogik/helios/router/http/health"

	authRoutes "gitlab.com/firelogik/helios/router/http/auth"
	mediaRoutes "gitlab.com/firelogik/helios/router/http/media"
	permissionRoutes "gitlab.com/firelogik/helios/router/http/permission"
	roleRoutes "gitlab.com/firelogik/helios/router/http/role"
	userRoutes "gitlab.com/firelogik/helios/router/http/user"

	logRoutes "gitlab.com/firelogik/helios/router/http/log"
	settingRoutes "gitlab.com/firelogik/helios/router/http/setting"

	swaggerRoutes "gitlab.com/firelogik/helios/router/http/swagger"
	websocketRoutes "gitlab.com/firelogik/helios/router/websocket"

	"gitlab.com/firelogik/helios/router/middleware"

	response "gitlab.com/firelogik/helios/utils/http"
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
