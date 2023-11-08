package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"html/template"
	"os"

	"gitlab.com/firelogik/helios/domain/model"
	"gopkg.in/gomail.v2"
)

type emailData struct {
	Data        interface{}
	FrontendURL string
}

type Emailer struct {
	Host           string
	Port           int
	SenderEmail    string
	SenderPassword string
	Template       *template.Template
}

type EmailContent struct {
	Body       string
	Attachment *string
	Data       interface{}
}

func (e Emailer) Notify(title string, content interface{}, recipient *model.User) error {
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

func (e Emailer) Subscribe(string, string) error { return errors.New("unsupported") }

func (e Emailer) Unsubscribe(string, string) error { return errors.New("unsupported") }
