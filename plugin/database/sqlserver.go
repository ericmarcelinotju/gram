package database

import (
	"fmt"

	"github.com/ericmarcelinotju/gram/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectSqlserver(configuration *config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s/%s?database=%s",
		configuration.User,
		configuration.Password,
		configuration.Host,
		configuration.Port,
		configuration.Instance,
		configuration.DB,
	)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})

	if err != nil {
		return nil, err
	}

	return db, nil
}
