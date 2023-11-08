package storage

import (
	"context"
	"fmt"

	"gitlab.com/firelogik/helios/config"
	"gitlab.com/firelogik/helios/constant"
	"gitlab.com/firelogik/helios/domain/module/setting"
	"gitlab.com/firelogik/helios/library/file"
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

func GetSFTPConfig(settingRepo setting.Repository) (*config.Storage, error) {
	ctx := context.Background()
	sftpHost, err := settingRepo.SelectSettingByName(ctx, constant.SFTPHost)
	if err != nil {
		return nil, err
	}
	sftpPort, err := settingRepo.SelectSettingByName(ctx, constant.SFTPPort)
	if err != nil {
		return nil, err
	}
	sftpUsername, err := settingRepo.SelectSettingByName(ctx, constant.SFTPUsername)
	if err != nil {
		return nil, err
	}
	sftpPassword, err := settingRepo.SelectSettingByName(ctx, constant.SFTPPassword)
	if err != nil {
		return nil, err
	}
	recordingFolder, err := settingRepo.SelectSettingByName(ctx, constant.SFTPStorageFolder)
	if err != nil {
		return nil, err
	}

	return &config.Storage{
		Path:     recordingFolder,
		Host:     fmt.Sprintf("%s:%s", sftpHost, sftpPort),
		Username: sftpUsername,
		Password: sftpPassword,
	}, nil
}
