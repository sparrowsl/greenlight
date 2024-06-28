package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFs embed.FS

type Mailer struct {
	dialer *mail.Client
	sender string
}

func New(host string, port int, username string, password string, sender string) Mailer {
	dialer, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(time.Second*5),
	)
	if err != nil {
		fmt.Println(err.Error())
		return Mailer{}
	}

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient string, templateFile string, data any) error {
	templ, err := template.New("email").ParseFS(templateFs, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = templ.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = templ.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = templ.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()
	msg.To(recipient)
	msg.From(m.sender)
	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.SetBodyHTMLTemplate(templ, data)

	// try to send the mail 3 times before finally aborting
	// if email successful, then return and cancel the other retries
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)
		if nil == err {
			return nil // if everything works, return successful/nil
		}

		// sleep 500ms between sending emails
		time.Sleep(time.Millisecond * 500)
	}

	return err
}
