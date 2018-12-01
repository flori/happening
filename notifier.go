package happening

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Notifier interface {
	Alert(check Check)
}

func NewNotifier(config ServerConfig) Notifier {
	kind := strings.ToLower(config.NOTIFIER_KIND)
	switch kind {
	case "mailcommand":
		return NewMailCommandNotifier(config)
	case "sendgrid":
		return NewSendgridNotifier(config)
	default:
		log.Panicf("unknown notifier of kind %s", config.NOTIFIER_KIND)
	}
	return nil
}

func mailSubject(environmentVariable string) string {
	return fmt.Sprintf(
		"Happening on %s has unhealthy checks",
		env12Factor(environmentVariable),
	)
}

func env12Factor(environmentVariable string) string {
	railsEnv, ok := os.LookupEnv(environmentVariable)
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
