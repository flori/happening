package happening

import (
	"log"
	"strings"
)

type Notifier interface {
	Alert(check Check)
	Resolve(check Check)
}

func NewNotifier(config ServerConfig) Notifier {
	kind := strings.ToLower(config.NOTIFIER_KIND)
	switch kind {
	case "":
		return NewNullNotifier()
	case "mailcommand":
		return NewMailCommandNotifier(config)
	case "sendgrid":
		return NewSendgridNotifier(config)
	default:
		log.Panicf("unknown notifier of kind %s", config.NOTIFIER_KIND)
	}
	return nil
}
