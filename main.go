package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/ericmarcelinotju/gram/command"
	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/library/email"

	"github.com/ericmarcelinotju/gram/repository/cache"
	"github.com/ericmarcelinotju/gram/repository/database"
	"github.com/ericmarcelinotju/gram/repository/job"
	"github.com/ericmarcelinotju/gram/repository/notifier"
	"github.com/ericmarcelinotju/gram/repository/storage"

	mediaStore "github.com/ericmarcelinotju/gram/repository/media"

	websocketStore "github.com/ericmarcelinotju/gram/repository/websocket"
	authModule "github.com/ericmarcelinotju/gram/module/auth"
	permissionModule "github.com/ericmarcelinotju/gram/module/permission"
	roleModule "github.com/ericmarcelinotju/gram/module/role"
	settingModule "github.com/ericmarcelinotju/gram/module/setting"
	userModule "github.com/ericmarcelinotju/gram/module/user"

	router "github.com/ericmarcelinotju/gram/router"

	"github.com/ericmarcelinotju/gram/domain/media"

	wsDomain "github.com/ericmarcelinotju/gram/domain/websocket"
)

// @securityDefinitions.apikey Auth
// @in header
// @name Authorization
func main() {
	// get configuration stucts via .env file
	configuration := config.NewConfig()

	// establish DB connection
	db, err := database.Connect(configuration.Database)
	if err != nil {
		log.Fatalln("[DATABASE] : ", err)
	}

	// establish cache connection
	redisCache, err := cache.ConnectRedis(configuration.Cache)
	if err != nil {
		log.Fatalln("[REDIS] : ", err)
	}

	// establish backup queue using redis
	backupQueue, err := job.ConnectQueue(
		&config.Queue{
			Name:            "queue",
			Number:          3,
			PrefetchLimit:   3,
			PollDuration:    time.Second * 3,
			ReportBatchSize: 3,
		},
		redisCache.Client(),
	)
	if err != nil {
		log.Fatalln("[BACKUP QUEUE] : ", err)
	}

	// initialize websocket dispatcher
	dispatcher, err := websocketStore.Init()
	if err != nil {
		log.Fatalln("[WEBSOCKET] : ", err)
	}

	var mediaStorage storage.Storage
	if configuration.MediaStorage != nil {
		// initialize File Manager
		mediaStorage, err = storage.InitFile(configuration.MediaStorage)
		if err != nil {
			log.Println("[FILE STORAGE] : ", err)
		}
	}

	// initialize scheduler for backup worker
	var firstBackupScheduler *job.Scheduler
	var secondBackupScheduler *job.Scheduler

	var forgotEmail *email.Emailer

	settingRepo := settingModule.NewRepository(db, redisCache)

	// Setup smtp from setting
	smtpConf, err := notifier.GetSMTPConfig(settingRepo)
	if err != nil {
		log.Println("[NOREPLY EMAIL] : ", err)
	}
	if smtpConf != nil {
		forgotEmail, err = notifier.InitEmailer(
			smtpConf,
			template.Must(template.ParseFiles("./email/template/forgot.html")),
		)
		if err != nil {
			log.Println("[FORGOT EMAIL] : ", err)
		}
	}

	authRepo := authModule.NewRepository(db, redisCache, forgotEmail)
	mediaRepo := mediaStore.New(mediaStorage)

	userRepo := userModule.NewRepository(db, mediaStorage)
	roleRepo := roleModule.NewRepository(db)
	permissionRepo := permissionModule.NewRepository(db)

	authSvc := authModule.NewService(authRepo, userRepo)
	mediaSvc := media.NewService(mediaRepo)

	userSvc := userModule.NewService(userRepo)
	roleSvc := roleModule.NewService(roleRepo)
	permissionSvc := permissionModule.NewService(permissionRepo)

	settingSvc := settingModule.NewService(settingRepo, firstBackupScheduler, secondBackupScheduler, forgotEmail)

	//websocket
	wsRepo, err := websocketStore.New(dispatcher)
	if err != nil {
		log.Fatalln("[WEBSOCKET] : ", err)
	}
	websocketSvc := wsDomain.NewService(wsRepo)

	command.ProcessCommands(db)

	router := router.NewHTTPHandler(
		authSvc,
		mediaSvc,

		userSvc,
		roleSvc,
		permissionSvc,

		settingSvc,

		websocketSvc,

		backupQueue.Connection,

		backupQueue,
	)
	log.Println("Start Listening to : " + configuration.Port)
	err = http.ListenAndServe(":"+configuration.Port, router)
	if err != nil {
		log.Fatalln(err)
	}
}
