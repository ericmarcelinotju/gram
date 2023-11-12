package setting

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/constant"
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/library/email"
	"github.com/ericmarcelinotju/gram/repository/job"
)

// Service defines Setting service behavior.
type Service interface {
	Read(context.Context) ([]dto.SettingDto, error)
	ReadByName(context.Context, string) (string, error)
	Save(context.Context, *dto.PostSettingDto) error

	GetSchedulerTime(context.Context, string) (int, int, error)
}

type service struct {
	repo                  Repository
	firstBackupScheduler  *job.Scheduler
	secondBackupScheduler *job.Scheduler
	forgotEmail           *email.Emailer
}

// NewService creates a new service struct
func NewService(
	repo Repository,
	firstBackupScheduler *job.Scheduler,
	secondBackupScheduler *job.Scheduler,
	forgotEmail *email.Emailer,
) *service {
	return &service{
		repo:                  repo,
		firstBackupScheduler:  firstBackupScheduler,
		secondBackupScheduler: secondBackupScheduler,
		forgotEmail:           forgotEmail,
	}
}

func (svc *service) Read(ctx context.Context) ([]dto.SettingDto, error) {
	return svc.repo.Select(ctx)
}

func (svc *service) ReadByName(ctx context.Context, setting string) (string, error) {
	return svc.repo.SelectByName(ctx, setting)
}

func (svc *service) Save(ctx context.Context, payload *dto.PostSettingDto) error {
	err := svc.repo.Save(ctx, payload.Name, payload.Value)
	if err != nil {
		return err
	}

	if payload.Name == constant.SMTPHost {
		svc.forgotEmail.Host = payload.Value
		return nil
	}
	if payload.Name == constant.SMTPPort {
		port, err := strconv.Atoi(payload.Value)
		if err != nil {
			return err
		}
		svc.forgotEmail.Port = port
		return nil
	}
	if payload.Name == constant.SMTPEmail {
		svc.forgotEmail.SenderEmail = payload.Value
		return nil
	}
	if payload.Name == constant.SMTPPassword {
		svc.forgotEmail.SenderPassword = payload.Value
		return nil
	}
	return nil
}

func (svc *service) GetSchedulerTime(ctx context.Context, settingKey string) (int, int, error) {
	backupTimeStr, err := svc.repo.SelectByName(ctx, settingKey)
	if err != nil {
		return 0, 0, err
	}
	backupTimeStrs := strings.Split(backupTimeStr, ":")
	if len(backupTimeStrs) != 2 {
		return 0, 0, errors.New("invalid backup time")
	}
	backupHour, err := strconv.Atoi(backupTimeStrs[0])
	if err != nil {
		return 0, 0, err
	}
	backupMinute, err := strconv.Atoi(backupTimeStrs[1])
	if err != nil {
		return 0, 0, err
	}
	return backupHour, backupMinute, nil
}

func (svc *service) GetSFTPConfig(ctx context.Context) (*config.Storage, error) {
	sftpHost, err := svc.repo.SelectByName(ctx, constant.SFTPHost)
	if err != nil {
		return nil, err
	}
	sftpPort, err := svc.repo.SelectByName(ctx, constant.SFTPPort)
	if err != nil {
		return nil, err
	}
	sftpUsername, err := svc.repo.SelectByName(ctx, constant.SFTPUsername)
	if err != nil {
		return nil, err
	}
	sftpPassword, err := svc.repo.SelectByName(ctx, constant.SFTPPassword)
	if err != nil {
		return nil, err
	}
	recordingFolder, err := svc.repo.SelectByName(ctx, constant.SFTPStorageFolder)
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
