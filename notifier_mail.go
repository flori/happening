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
			`Happening on %s has healthy check "%s", problem was resolved`,
			mail.env12Factor(),
			mail.Check.Name,
		)
	} else {
		return fmt.Sprintf(
			`Happening on %s has unhealthy check "%s"`,
			mail.env12Factor(),
			mail.Check.Name,
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
	if mail.DrilldownURL != "" {
		text += fmt.Sprintf(
			"\n\nDrill down via this URL: %s/search/name:%s",
			mail.DrilldownURL,
			mail.Check.Name,
		)
	}
	return text
}
