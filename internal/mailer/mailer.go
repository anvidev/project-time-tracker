package mailer

import "embed"

//go:embed "templates"
var templates embed.FS

var (
	NotifyEmptyDay = "notify_empty_day.html"
)

type Mailer interface {
	Send(to []string, subject, tmpl string, data any) error
}
