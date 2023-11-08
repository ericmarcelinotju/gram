package notifier

import (
	"context"
	"html/template"
	"strconv"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/constant"
	"github.com/ericmarcelinotju/gram/domain/module/setting"
	"github.com/ericmarcelinotju/gram/library/email"
)

func InitEmailer(configuration *config.Email, template *template.Template) (*email.Emailer, error) {
	port, err := strconv.Atoi(configuration.Port)
	if err != nil {
		return nil, err
	}
	return &email.Emailer{
		Host:           configuration.Host,
		Port:           port,
		SenderEmail:    configuration.Email,
		SenderPassword: configuration.Password,
		Template:       template,
	}, nil
}

func GetSMTPConfig(settingRepo setting.Repository) (*config.Email, error) {
	ctx := context.Background()
	noreplyHost, err := settingRepo.SelectSettingByName(ctx, constant.SMTPHost)
	if err != nil {
		return nil, err
	}
	noreplyPort, err := settingRepo.SelectSettingByName(ctx, constant.SMTPPort)
	if err != nil {
		return nil, err
	}
	noreplyEmail, err := settingRepo.SelectSettingByName(ctx, constant.SMTPEmail)
	if err != nil {
		return nil, err
	}
	noreplyPassword, err := settingRepo.SelectSettingByName(ctx, constant.SMTPPassword)
	if err != nil {
		return nil, err
	}
	return &config.Email{
		Host:     noreplyHost,
		Port:     noreplyPort,
		Email:    noreplyEmail,
		Password: noreplyPassword,
	}, nil
}
