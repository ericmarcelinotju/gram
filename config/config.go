package config

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ericmarcelinotju/gram/utils/env"
	"github.com/ericmarcelinotju/gram/utils/job"
	"github.com/gin-gonic/gin"
)

// Config is a struct that contains configuration variables
type Config struct {
	Version       string
	Environment   string
	Host          url.URL
	Port          string
	Net           *Net
	Database      *Database
	DatabaseMeter *Database
	Queue         *Queue
	Cache         *Cache
	Secret        string
	MediaStorage  *Storage
}

// Net is a struct that contains net client's configuration variables
type Net struct {
	Timeout time.Duration
}

// Database is a struct that contains DB's configuration variables
type Database struct {
	Driver   string
	Host     string
	Port     string
	Instance string
	User     string
	DB       string
	Password string
}

// Cache is a struct that contains cache's configuration variables
type Cache struct {
	Driver        string
	Password      string
	Host          string
	Port          string
	DefaultExpiry time.Duration
}

// Queue is a struct that contains consumer's configuration variables
type Queue struct {
	Name            string
	Number          int
	PrefetchLimit   int64
	PollDuration    time.Duration
	ReportBatchSize int
}

// Storage is a struct that contains Storage's configuration variables
type Storage struct {
	Path string
	URL  string

	Host     string
	Username string
	Password string
}

// Email is a struct that contains email's configuration variables
type Email struct {
	Host     string
	Port     string
	Email    string
	Password string
}

func setLogOutputFile(filename string) {
	var file *os.File
	file, err := os.OpenFile(filename+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		file, _ = os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	}
	log.SetOutput(file)
}

func HandleLogFile() {
	go func() {
		currTimeStr := time.Now().UTC().Format("2006-01-02")
		setLogOutputFile(currTimeStr)

		jobTicker := &job.JobTicker{}
		jobTicker.UpdateTimer()
		for {
			<-jobTicker.Timer.C

			currTimeStr := time.Now().UTC().Format("2006-01-02")

			setLogOutputFile(currTimeStr)

			jobTicker.UpdateTimer()
		}
	}()
}

var Instance *Config

// NewConfig creates a new Config struct
func NewConfig(envFileName string) *Config {
	envFile = envFileName
	env.CheckDotEnv(envFileName)
	environment := env.MustGet("ENV")
	if environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
		HandleLogFile()
	}

	port := env.MustGet("PORT")
	if port == "" {
		port = "3000"
	}

	hostStr := env.MustGet("HOST")
	host, err := url.Parse(hostStr)
	if host == nil || err != nil {
		panic("Error when parsing host : '" + hostStr + "'")
	}

	cacheExpiryInt, err := strconv.Atoi(env.MustGet("CACHE_EXPIRY"))
	if err != nil {
		panic("Error when parsing cache expiry")
	}
	cacheExpiry := time.Millisecond * time.Duration(cacheExpiryInt)

	netTimeoutInt, err := strconv.Atoi(env.MustGet("NET_TIMEOUT"))
	if err != nil {
		panic("Error when parsing cache expiry")
	}
	netTimeout := time.Millisecond * time.Duration(netTimeoutInt)

	config := &Config{
		Version:     "1.0.0",
		Environment: environment,
		Host:        *host,
		Port:        port,
		Net: &Net{
			Timeout: netTimeout,
		},
		Database: &Database{
			Driver:   env.MustGet("DB_DRIVER"),
			Host:     env.MustGet("DB_HOST"),
			Port:     env.MustGet("DB_PORT"),
			Instance: env.Get("DB_INSTANCE"),
			User:     env.MustGet("DB_USER"),
			DB:       env.MustGet("DB_NAME"),
			Password: env.Get("DB_PASSWORD"),
		},
		Cache: &Cache{
			Driver:        env.MustGet("CACHE_DRIVER"),
			Password:      env.Get("CACHE_PASSWORD"),
			Host:          env.MustGet("CACHE_HOST"),
			Port:          env.MustGet("CACHE_PORT"),
			DefaultExpiry: cacheExpiry,
		},
	}

	mediaPath := env.Get("MEDIA_PATH")
	if mediaPath != "" {
		config.MediaStorage = &Storage{
			Path: env.MustGet("MEDIA_PATH"),
			URL:  env.MustGet("MEDIA_URL"),
		}
	}

	Instance = config
	return config
}

var configInstance *Config
var once sync.Once
var envFile = ".env"

func Get() *Config {
	once.Do(func() {
		configInstance = NewConfig(envFile)
	})

	return configInstance
}
