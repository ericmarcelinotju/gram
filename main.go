package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/ericmarcelinotju/gram/command"
	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/library/email"

	"github.com/ericmarcelinotju/gram/data/cache"
	"github.com/ericmarcelinotju/gram/data/database"
	"github.com/ericmarcelinotju/gram/data/job"
	"github.com/ericmarcelinotju/gram/data/notifier"
	"github.com/ericmarcelinotju/gram/data/storage"

	mediaStore "github.com/ericmarcelinotju/gram/data/media"
	authStore "github.com/ericmarcelinotju/gram/data/module/auth"

	logStore "github.com/ericmarcelinotju/gram/data/module/log"
	permissionStore "github.com/ericmarcelinotju/gram/data/module/permission"
	roleStore "github.com/ericmarcelinotju/gram/data/module/role"
	settingStore "github.com/ericmarcelinotju/gram/data/module/setting"
	userStore "github.com/ericmarcelinotju/gram/data/module/user"
	websocketStore "github.com/ericmarcelinotju/gram/data/websocket"

	router "github.com/ericmarcelinotju/gram/router"

	"github.com/ericmarcelinotju/gram/domain/media"
	"github.com/ericmarcelinotju/gram/domain/module/auth"
	logService "github.com/ericmarcelinotju/gram/domain/module/log"
	"github.com/ericmarcelinotju/gram/domain/module/permission"
	"github.com/ericmarcelinotju/gram/domain/module/role"
	"github.com/ericmarcelinotju/gram/domain/module/setting"
	"github.com/ericmarcelinotju/gram/domain/module/user"

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

	settingRepo := settingStore.New(db, redisCache)

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

	authRepo := authStore.New(db, redisCache, forgotEmail)
	mediaRepo := mediaStore.New(mediaStorage)

	userRepo := userStore.New(db, mediaStorage)
	roleRepo := roleStore.New(db)
	permissionRepo := permissionStore.New(db)

	logRepo := logStore.New(db)

	authSvc := auth.NewService(authRepo, userRepo, logRepo)
	mediaSvc := media.NewService(mediaRepo)

	userSvc := user.NewService(userRepo, roleRepo)
	roleSvc := role.NewService(roleRepo, permissionRepo)
	permissionSvc := permission.NewService(permissionRepo)

	logSvc := logService.NewService(logRepo)
	settingSvc := setting.NewService(settingRepo, firstBackupScheduler, secondBackupScheduler, forgotEmail)

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

		logSvc,
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
