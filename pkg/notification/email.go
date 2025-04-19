package notification

import (
	"bytes"
	"embed"
	"github.com/wneessen/go-mail"
	ht "html/template"
	tt "text/template"
	"time"
)

//go:embed "templates"
var templateFS embed.FS

type Email interface {
	SendEmail(recipient, templateFile string, data any) error
}

type email struct {
	client *mail.Client
	sender string
}

func (e *email) SendEmail(recipient, templateFile string, data any) error {
	textTmpl, err := tt.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlTmpl, err := ht.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = htmlTmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()

	err = msg.To(recipient)
	if err != nil {
		return err
	}

	err = msg.From(e.sender)
	if err != nil {
		return err
	}

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	for i := 1; i <= 3; i++ {
		err = e.client.DialAndSend(msg)
		if err == nil {
			return nil
		}
		if i != 3 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return err
}

func NewEmail(host string, port int, username, password, sender string) (Email, error) {
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	e := &email{
		client: client,
		sender: sender,
	}

	return e, nil
}
