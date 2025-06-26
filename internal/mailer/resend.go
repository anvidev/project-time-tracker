package mailer

import (
	"bytes"
	"html/template"

	"github.com/resend/resend-go/v2"
)

type resendMailer struct {
	client *resend.Client
	from   string // format: "name <email>"
}

func NewResendMailer(apikey, from string) Mailer {
	return &resendMailer{
		from:   from,
		client: resend.NewClient(apikey),
	}
}

func (m resendMailer) Send(to []string, subject, tmpl string, data any) error {
	html, err := template.ParseFS(templates, "templates/"+tmpl)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := html.ExecuteTemplate(body, "body", data); err != nil {
		return err
	}

	email := &resend.SendEmailRequest{
		From:    m.from,
		To:      to,
		Subject: subject,
		Html:    body.String(),
	}

	_, err = m.client.Emails.Send(email)

	return err
}
