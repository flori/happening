package happening

import (
	"fmt"
	"os"
)

type NotifierMail struct {
	Check               Check
	EnvironmentVariable string
	DrilldownURL        string
}

func (mail NotifierMail) Subject() string {
	return fmt.Sprintf(
		"Happening on %s has unhealthy checks",
		mail.env12Factor(),
	)
}

func (mail NotifierMail) env12Factor() string {
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

func (mail NotifierMail) Text() string {
	text := mail.Check.String()
	if mail.DrilldownURL != "" {
		text += fmt.Sprintf(
			"\n\nDrill down via this URL: %s/search/name:%s",
			mail.DrilldownURL,
			mail.Check.Name,
		)
	}
	return text
}
