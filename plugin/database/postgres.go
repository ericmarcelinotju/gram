package database

import (
	"fmt"

	"github.com/ericmarcelinotju/gram/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(configuration *config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		configuration.Host,
		configuration.User,
		configuration.Password,
		configuration.DB,
		configuration.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})

	if err != nil {
		return nil, err
	}

	return db, nil
}
