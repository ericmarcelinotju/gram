package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"gitlab.com/firelogik/helios/command"
	"gitlab.com/firelogik/helios/config"
	"gitlab.com/firelogik/helios/library/email"

	"gitlab.com/firelogik/helios/data/cache"
	"gitlab.com/firelogik/helios/data/database"
	"gitlab.com/firelogik/helios/data/job"
	"gitlab.com/firelogik/helios/data/notifier"
	"gitlab.com/firelogik/helios/data/storage"

	mediaStore "gitlab.com/firelogik/helios/data/media"
	authStore "gitlab.com/firelogik/helios/data/module/auth"

	logStore "gitlab.com/firelogik/helios/data/module/log"
	permissionStore "gitlab.com/firelogik/helios/data/module/permission"
	roleStore "gitlab.com/firelogik/helios/data/module/role"
	settingStore "gitlab.com/firelogik/helios/data/module/setting"
	userStore "gitlab.com/firelogik/helios/data/module/user"
	websocketStore "gitlab.com/firelogik/helios/data/websocket"

	router "gitlab.com/firelogik/helios/router"

	"gitlab.com/firelogik/helios/domain/media"
	"gitlab.com/firelogik/helios/domain/module/auth"
	logService "gitlab.com/firelogik/helios/domain/module/log"
	"gitlab.com/firelogik/helios/domain/module/permission"
	"gitlab.com/firelogik/helios/domain/module/role"
	"gitlab.com/firelogik/helios/domain/module/setting"
	"gitlab.com/firelogik/helios/domain/module/user"

	wsDomain "gitlab.com/firelogik/helios/domain/websocket"
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
