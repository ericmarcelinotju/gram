package database

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/utils/env"
	"github.com/glebarez/sqlite"
	sqliteGo "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

func MD5(val string) string {
	if len(val) == 0 {
		return val
	}
	hash := md5.Sum([]byte(val))
	return hex.EncodeToString(hash[:])
}

var once sync.Once

func ConnectSqlite(configuration *config.Database) (*gorm.DB, error) {
	const CustomDriverName = "sqlite3_extended"
	var File = env.GetRootPath(configuration.DB)
	//dsn := fmt.Sprintf(
	//	"file:%s.db?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha256",
	//	File,
	//	configuration.User,
	//	configuration.Password,
	//)
	once.Do(func() {
		sql.Register(CustomDriverName,
			&sqliteGo.SQLiteDriver{
				ConnectHook: func(conn *sqliteGo.SQLiteConn) error {
					err := conn.RegisterFunc(
						"uuid_generate_v4",
						func(arguments ...interface{}) (string, error) {
							return uuid.NewV4().String(), nil // Return a string value.
						},
						true,
					)
					return err
				},
			},
		)
	})

	conn, err := sql.Open(CustomDriverName, File)
	if err != nil {
		panic(err)
	}
	return gorm.Open(sqlite.Dialector{
		DriverName: CustomDriverName,
		DSN:        File,
		Conn:       conn,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	//

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

	//return gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
}
