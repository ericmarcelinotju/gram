package database

import (
	"errors"

	"gitlab.com/firelogik/helios/config"
	"gorm.io/gorm"
)

func Connect(configuration *config.Database) (*gorm.DB, error) {
	if configuration.Driver == "sqlite" {
		return ConnectSqlite(configuration)
	} else if configuration.Driver == "postgres" {
		return ConnectPostgres(configuration)
	} else if configuration.Driver == "mysql" {
		return ConnectMysql(configuration)
	} else if configuration.Driver == "sqlserver" {
		return ConnectSqlserver(configuration)
	} else {
		return nil, errors.New("database driver unsupported")
	}
}
