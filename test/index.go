package test

import (
	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/data/cache"
	"github.com/ericmarcelinotju/gram/data/database"
	"github.com/ericmarcelinotju/gram/data/notifier"
	"github.com/ericmarcelinotju/gram/data/storage"
	"github.com/gin-gonic/gin"

	authStore "github.com/ericmarcelinotju/gram/data/module/auth"
	logStore "github.com/ericmarcelinotju/gram/data/module/log"
	permissionStore "github.com/ericmarcelinotju/gram/data/module/permission"
	roleStore "github.com/ericmarcelinotju/gram/data/module/role"
	settingStore "github.com/ericmarcelinotju/gram/data/module/setting"
	userStore "github.com/ericmarcelinotju/gram/data/module/user"

	"github.com/ericmarcelinotju/gram/domain/module/auth"
	"github.com/ericmarcelinotju/gram/domain/module/log"
	"github.com/ericmarcelinotju/gram/domain/module/permission"
	"github.com/ericmarcelinotju/gram/domain/module/role"
	"github.com/ericmarcelinotju/gram/domain/module/setting"
	"github.com/ericmarcelinotju/gram/domain/module/user"
)

var AUTHREPO auth.Repository
var PERMISSIONREPO permission.Repository
var ROLEREPO role.Repository
var USERREPO user.Repository
var LOGREPO log.Repository
var SETTINGREPO setting.Repository

func initRepositories() {
	// get configuration stucts via .env file
	configuration := config.NewConfig()

	// establish DB connection
	db, _ := database.Connect(configuration.Database)

	// establish cache connection
	redisCache, _ := cache.ConnectRedis(configuration.Cache)

	var fileStorage storage.Storage
	if configuration.MediaStorage != nil {
		// initialize File Manager
		fileStorage, _ = storage.InitFile(configuration.MediaStorage)
	}

	// initialize Forgot Password Email
	var forgotEmail notifier.Notifier

	AUTHREPO = authStore.New(db, redisCache, forgotEmail)

	PERMISSIONREPO = permissionStore.New(db)
	ROLEREPO = roleStore.New(db)
	USERREPO = userStore.New(db, fileStorage)

	LOGREPO = logStore.New(db)
	SETTINGREPO = settingStore.New(db, redisCache)
}

var AUTHSVC auth.Service
var PERMISSIONSVC permission.Service
var ROLESVC role.Service
var USERSVC user.Service
var LOGSVC log.Service
var SETTINGSVC setting.Service

func initServices() {
	AUTHSVC = auth.NewService(AUTHREPO, USERREPO, LOGREPO)
	PERMISSIONSVC = permission.NewService(PERMISSIONREPO)
	ROLESVC = role.NewService(ROLEREPO, PERMISSIONREPO)
	USERSVC = user.NewService(USERREPO, ROLEREPO)

	LOGSVC = log.NewService(LOGREPO)
	SETTINGSVC = setting.NewService(SETTINGREPO, nil, nil, nil)

	// migrate := command.MigrateCommandFactory([]domain.S{PERMISSIONSVC, ROLESVC, USERSVC, LOGSVC, SETTINGSVC})

	// migrate(context.Background())

	// seed := command.SeedingCommandFactory([]domain.SeederService{PERMISSIONSVC, ROLESVC, SETTINGSVC, USERSVC})

	// seed(context.Background())
}

func init() {
	initRepositories()
	initServices()
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	return router
}
