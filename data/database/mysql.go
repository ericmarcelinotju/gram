package database

import (
	"fmt"

	"github.com/ericmarcelinotju/gram/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectMysql(configuration *config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		configuration.User,
		configuration.Password,
		configuration.Host,
		configuration.Port,
		configuration.DB,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})

	if err != nil {
		return nil, err
	}

	return db, nil
}
