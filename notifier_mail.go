package happening

import (
	"fmt"
	"os"
)

type NotifierMail struct {
	Check               Check
	EnvironmentVariable string
	DrilldownURL        string
	Resolved            bool
}

func (mail NotifierMail) Subject() string {
	if mail.Resolved {
		return fmt.Sprintf(
			`Resolved: Happening on %s check "%s": %s`,
			mail.env12Factor(),
			mail.Check.Name,
			mail.Check.State(),
		)
	} else {
		return fmt.Sprintf(
			`Problem: Happening on %s check "%s": %s`,
			mail.env12Factor(),
			mail.Check.Name,
			mail.Check.State(),
		)
	}
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
	if mail.DrilldownURL != "" && mail.Check.LastEventId != nil {
		text += fmt.Sprintf(
			"\n\nDrill down via this URL: %s/search/id:%s",
			mail.DrilldownURL,
			*mail.Check.LastEventId,
		)
	}
	return text
}
