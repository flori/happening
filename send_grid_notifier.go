package happening

import (
	"log"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridNotifier struct {
	ENVIRONMENT_VARIABLE string
	SENDGRID_API_KEY     string
	NO_REPLY_NAME        string
	NO_REPLY_EMAIL       string
	CONTACT_NAME         string
	CONTACT_EMAIL        string
}

func NewSendgridNotifier(config ServerConfig) Notifier {
	return &SendGridNotifier{
		ENVIRONMENT_VARIABLE: config.NOTIFIER_ENVIRONMENT_VARIABLE,
		SENDGRID_API_KEY:     config.NOTIFIER_SENDGRID_API_KEY,
		NO_REPLY_NAME:        config.NOTIFIER_NO_REPLY_NAME,
		NO_REPLY_EMAIL:       config.NOTIFIER_NO_REPLY_EMAIL,
		CONTACT_NAME:         config.NOTIFIER_CONTACT_NAME,
		CONTACT_EMAIL:        config.NOTIFIER_CONTACT_EMAIL,
	}
}

func (notifier *SendGridNotifier) buildMail(text string) []byte {
	from := mail.NewEmail(notifier.NO_REPLY_NAME, notifier.NO_REPLY_EMAIL)
	to := mail.NewEmail(notifier.CONTACT_NAME, notifier.CONTACT_EMAIL)
	replyTo := mail.NewEmail(notifier.CONTACT_NAME, notifier.CONTACT_EMAIL)
	subject := mailSubject(notifier.ENVIRONMENT_VARIABLE)
	content := mail.NewContent("text/plain", text)
	m := mail.NewV3MailInit(from, subject, to, content)
	m.SetReplyTo(replyTo)
	return mail.GetRequestBody(m)
}

func (notifier *SendGridNotifier) sendMail(text string) {
	if notifier.SENDGRID_API_KEY == "" {
		log.Panicln("Sendgrid API key required in environment configuration")
	}
	request := sendgrid.GetRequest(notifier.SENDGRID_API_KEY, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = notifier.buildMail(text)
	if _, err := sendgrid.API(request); err != nil {
		log.Panic(err)
	}
}

func (notifier *SendGridNotifier) Alert(check Check) {
	go notifier.sendMail(check.String())
}
