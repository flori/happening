package happening

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type MailCommandNotifier struct {
	EnvironmentVariable string
	DrilldownURL        string
	MailCommand         string
	ContactName         string
	ContactEmail        string
}

func NewMailCommandNotifier(config ServerConfig) Notifier {
	return &MailCommandNotifier{
		EnvironmentVariable: config.NOTIFIER_ENVIRONMENT_VARIABLE,
		DrilldownURL:        config.NOTIFIER_DRILLDOWN_URL,
		MailCommand:         config.NOTIFIER_MAIL_COMMAND,
		ContactName:         config.NOTIFIER_CONTACT_NAME,
		ContactEmail:        config.NOTIFIER_CONTACT_EMAIL,
	}
}

func (notifier *MailCommandNotifier) sendMail(notifierMail NotifierMail) {
	to := fmt.Sprintf("%s <%s>", notifier.ContactName, notifier.ContactEmail)
	path, err := exec.LookPath(notifier.MailCommand)
	if err != nil {
		log.Printf("error: %v", err)
	}
	cmd := exec.Command(
		path,
		"-s",
		notifierMail.Subject(),
		to,
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf(`REPLYTO="%s"`, to))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, notifierMail.Text())
	}()
	if err = cmd.Run(); err != nil {
		log.Panic(err)
	}
}

func (notifier *MailCommandNotifier) Alert(check Check) {
	go notifier.sendMail(
		NotifierMail{
			Check:               check,
			EnvironmentVariable: notifier.EnvironmentVariable,
			DrilldownURL:        notifier.DrilldownURL,
		},
	)
}
