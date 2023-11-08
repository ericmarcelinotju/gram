package database

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/glebarez/sqlite"
	"gitlab.com/firelogik/helios/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MD5(val string) string {
	if len(val) == 0 {
		return val
	}
	hash := md5.Sum([]byte(val))
	return hex.EncodeToString(hash[:])
}

func ConnectSqlite(configuration *config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"file:%s.db?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha256",
		configuration.DB,
		configuration.User,
		configuration.Password,
	)

	// sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
	// 	ConnectHook: func(conn *sqlite3.SQLiteConn) error {
	// 		if err := conn.RegisterFunc("MD5", MD5, false); err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	},
	// })

	// conn, err := sql.Open("sqlite3", dsn)
	// if err != nil {
	// 	log.Fatal("Failed to open database:", err)
	// }

	return gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
}
