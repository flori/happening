package happening

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type MailCommandNotifier struct {
	ENVIRONMENT_VARIABLE string
	MAIL_COMMAND         string
	CONTACT_NAME         string
	CONTACT_EMAIL        string
}

func NewMailCommandNotifier(config ServerConfig) Notifier {
	return &MailCommandNotifier{
		ENVIRONMENT_VARIABLE: config.NOTIFIER_ENVIRONMENT_VARIABLE,
		MAIL_COMMAND:         config.NOTIFIER_MAIL_COMMAND,
		CONTACT_NAME:         config.NOTIFIER_CONTACT_NAME,
		CONTACT_EMAIL:        config.NOTIFIER_CONTACT_EMAIL,
	}
}

func (notifier *MailCommandNotifier) sendMail(text string) {
	to := fmt.Sprintf("%s <%s>", notifier.CONTACT_NAME, notifier.CONTACT_EMAIL)
	path, err := exec.LookPath(notifier.MAIL_COMMAND)
	if err != nil {
		log.Printf("error: %v", err)
	}
	cmd := exec.Command(
		path,
		"-s",
		mailSubject(notifier.ENVIRONMENT_VARIABLE),
		to,
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf(`REPLYTO="%s"`, to))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, text)
	}()
	if err = cmd.Run(); err != nil {
		log.Panic(err)
	}
}

func (notifier *MailCommandNotifier) Alert(check Check) {
	go notifier.sendMail(check.String())
}
