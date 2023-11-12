package notifier

import (
	"context"
	"html/template"
	"strconv"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/constant"
	"github.com/ericmarcelinotju/gram/library/email"
	settingModule "github.com/ericmarcelinotju/gram/module/setting"
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

func GetSMTPConfig(settingRepo settingModule.Repository) (*config.Email, error) {
	ctx := context.Background()
	noreplyHost, err := settingRepo.SelectByName(ctx, constant.SMTPHost)
	if err != nil {
		return nil, err
	}
	noreplyPort, err := settingRepo.SelectByName(ctx, constant.SMTPPort)
	if err != nil {
		return nil, err
	}
	noreplyEmail, err := settingRepo.SelectByName(ctx, constant.SMTPEmail)
	if err != nil {
		return nil, err
	}
	noreplyPassword, err := settingRepo.SelectByName(ctx, constant.SMTPPassword)
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
