package storage

import (
	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/library/file"
)

func InitFile(configuration *config.Storage) (*file.File, error) {
	return file.NewFileManager(configuration.Path)
}

func InitFTP(configuration *config.Storage) (*file.FTP, error) {
	return file.NewFTPManager(
		configuration.Path,
		configuration.Host,
		configuration.Username,
		configuration.Password,
	)
}
