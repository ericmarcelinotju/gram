package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/ericmarcelinotju/gram/command"
	"github.com/ericmarcelinotju/gram/config"

	"github.com/ericmarcelinotju/gram/plugins/cache"
	"github.com/ericmarcelinotju/gram/plugins/database"
	"github.com/ericmarcelinotju/gram/plugins/job"
	"github.com/ericmarcelinotju/gram/plugins/notifier"
	"github.com/ericmarcelinotju/gram/plugins/storage"

	authModule "github.com/ericmarcelinotju/gram/module/auth"
	permissionModule "github.com/ericmarcelinotju/gram/module/permission"
	roleModule "github.com/ericmarcelinotju/gram/module/role"
	settingModule "github.com/ericmarcelinotju/gram/module/setting"
	userModule "github.com/ericmarcelinotju/gram/module/user"
	websocketStore "github.com/ericmarcelinotju/gram/plugins/websocket"

	exampleScheduler "github.com/ericmarcelinotju/gram/scheduler/example"

	router "github.com/ericmarcelinotju/gram/router"
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

	// establish job queue using redis
	jobQueue, err := job.ConnectQueue(
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
	dispatcher, err := websocketStore.NewDispatcher()
	if err != nil {
		log.Fatalln("[WEBSOCKET] : ", err)
	}

	var mediaStorage storage.Storage
	if configuration.MediaStorage != nil {
		// initialize File Manager
		mediaStorage, err = storage.NewFileStorage(configuration.MediaStorage)
		if err != nil {
			log.Println("[FILE STORAGE] : ", err)
		}
	}

	// initialize scheduler for backup worker
	var scheduler *job.Scheduler

	var forgotEmail *notifier.EmailNotifier

	settingRepo := settingModule.NewRepository(db, redisCache)

	authRepo := authModule.NewRepository(db, redisCache, forgotEmail)

	userRepo := userModule.NewRepository(db, mediaStorage, dispatcher)
	roleRepo := roleModule.NewRepository(db)
	permissionRepo := permissionModule.NewRepository(db)

	authSvc := authModule.NewService(authRepo, userRepo)

	userSvc := userModule.NewService(userRepo)
	roleSvc := roleModule.NewService(roleRepo)
	permissionSvc := permissionModule.NewService(permissionRepo)

	settingSvc := settingModule.NewService(settingRepo, scheduler, forgotEmail)

	// Setup smtp from setting
	smtpConf, err := settingSvc.GetSMTPConfig(context.Background())
	if err != nil {
		log.Println("[NOREPLY EMAIL] : ", err)
	}
	if smtpConf != nil {
		forgotEmail, err = notifier.NewEmailNotifier(
			smtpConf,
			template.Must(template.ParseFiles("./email/template/forgot.html")),
		)
		if err != nil {
			log.Println("[FORGOT EMAIL] : ", err)
		}
	}

	command.ProcessCommands(db)

	exampleScheduler, err := exampleScheduler.NewScheduler(jobQueue)
	if err != nil {
		log.Fatalln(err)
	}
	exampleScheduler.Start()

	router := router.NewHTTPHandler(
		authSvc,

		userSvc,
		roleSvc,
		permissionSvc,

		settingSvc,

		// TODO :: fix this shit
		jobQueue.Connection,
		jobQueue,
	)
	log.Println("Start Listening to : " + configuration.Port)
	err = http.ListenAndServe(":"+configuration.Port, router)
	if err != nil {
		log.Fatalln(err)
	}
}
