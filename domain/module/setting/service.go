package setting

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/ericmarcelinotju/gram/constant"
	"github.com/ericmarcelinotju/gram/data/job"
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/ericmarcelinotju/gram/library/email"
)

// Service defines Setting service behavior.
type Service interface {
	ReadSetting(context.Context) ([]model.Setting, error)
	ReadSettingByName(context.Context, string) (string, error)
	SaveSetting(context.Context, *model.Setting) error

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

func (svc *service) ReadSetting(ctx context.Context) ([]model.Setting, error) {
	return svc.repo.SelectSetting(ctx)
}

func (svc *service) ReadSettingByName(ctx context.Context, setting string) (string, error) {
	return svc.repo.SelectSettingByName(ctx, setting)
}

func (svc *service) SaveSetting(ctx context.Context, payload *model.Setting) error {
	err := svc.repo.SaveSetting(ctx, payload.Name, payload.Value)
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
	backupTimeStr, err := svc.repo.SelectSettingByName(ctx, settingKey)
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
