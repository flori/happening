package happening

import (
	"fmt"
	"os"
)

type EventNotifierMail struct {
	Event               Event
	EnvironmentVariable string
	DrilldownURL        string
}

func (mail EventNotifierMail) Subject() string {
	return fmt.Sprintf(
		`Event: Happening on %s event "%s": %s`,
		mail.env12Factor(),
		mail.Event.Name,
		mail.Event.Result(),
	)
}

func (mail EventNotifierMail) env12Factor() string {
	railsEnv, ok := os.LookupEnv(mail.EnvironmentVariable)
	if !ok {
		railsEnv = "development"
	}
	staging, ok := os.LookupEnv("STAGING")
	if !ok {
		staging = "0"
	}
	if railsEnv == "production" && staging == "1" {
		return "staging"
	} else {
		return railsEnv
	}
}

func (mail EventNotifierMail) Text() string {
	text := mail.Event.String()
	text += fmt.Sprintf(
		"\n\nDrill down via this URL: %s/search/id:%s",
		mail.DrilldownURL,
		mail.Event.Id,
	)
	return text
}
