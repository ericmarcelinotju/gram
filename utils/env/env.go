package env

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	local      = "local"
	staging    = "staging"
	production = "production"
	env        = "ENV"
)

// MustGet will return the env or panic if not present.
func MustGet(key string) string {
	val := os.Getenv(key)
	if val == "" && key != "PORT" {
		panic("Env key missing " + key)
	}
	return val
}

// Get will return the env or panic if not present.
func Get(key string) string {
	return os.Getenv(key)
}

// CheckDotEnv loads environment variables from .env file for development environment
func CheckDotEnv(envFileName string) {
	err := godotenv.Overload(GetRootPath(envFileName))
	if err != nil && os.Getenv(env) == local {
		log.Println("Error loading .env file")
	}

}

// GetRootPath returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
func GetRootPath(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}
