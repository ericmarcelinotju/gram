package test

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/config"
	"gitlab.com/firelogik/helios/data/cache"
	"gitlab.com/firelogik/helios/data/database"
	"gitlab.com/firelogik/helios/data/notifier"
	"gitlab.com/firelogik/helios/data/storage"

	authStore "gitlab.com/firelogik/helios/data/module/auth"
	logStore "gitlab.com/firelogik/helios/data/module/log"
	permissionStore "gitlab.com/firelogik/helios/data/module/permission"
	roleStore "gitlab.com/firelogik/helios/data/module/role"
	settingStore "gitlab.com/firelogik/helios/data/module/setting"
	userStore "gitlab.com/firelogik/helios/data/module/user"

	"gitlab.com/firelogik/helios/domain/module/auth"
	"gitlab.com/firelogik/helios/domain/module/log"
	"gitlab.com/firelogik/helios/domain/module/permission"
	"gitlab.com/firelogik/helios/domain/module/role"
	"gitlab.com/firelogik/helios/domain/module/setting"
	"gitlab.com/firelogik/helios/domain/module/user"
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
