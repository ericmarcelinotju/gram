package notifier

import (
	"bytes"
	"crypto/tls"
	"errors"
	"html/template"
	"os"
	"strconv"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/dto"
	"gopkg.in/gomail.v2"
)

type emailData struct {
	Data        interface{}
	FrontendURL string
}

type EmailContent struct {
	Body       string
	Attachment *string
	Data       interface{}
}

type EmailNotifier struct {
	Host           string
	Port           int
	SenderEmail    string
	SenderPassword string
	Template       *template.Template
}

func (e EmailNotifier) Notify(title string, content interface{}, recipient *dto.UserDto) error {
	// Sender data.
	from := e.SenderEmail
	password := e.SenderPassword

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetAddressHeader("To", recipient.Email, recipient.Username)
	m.SetHeader("Subject", title)

	emailContent, ok := content.(EmailContent)
	if !ok {
		return errors.New("incorrect email content")
	}

	if emailContent.Attachment != nil {
		m.Attach(*emailContent.Attachment)
	}

	emailData := emailData{
		Data:        emailContent.Data,
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}

	if e.Template != nil {
		var tpl bytes.Buffer
		e.Template.Execute(&tpl, emailData)

		m.SetBody("text/html", tpl.String())
	} else {
		m.SetBody("text/plain", emailContent.Body)
	}

	d := gomail.NewDialer(e.Host, e.Port, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (e EmailNotifier) Subscribe(string, string) error { return errors.New("unsupported") }

func (e EmailNotifier) Unsubscribe(string, string) error { return errors.New("unsupported") }

func NewEmailNotifier(configuration *config.Email, template *template.Template) (*EmailNotifier, error) {
	port, err := strconv.Atoi(configuration.Port)
	if err != nil {
		return nil, err
	}
	return &EmailNotifier{
		Host:           configuration.Host,
		Port:           port,
		SenderEmail:    configuration.Email,
		SenderPassword: configuration.Password,
		Template:       template,
	}, nil
}
