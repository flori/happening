package happening

import (
	"log"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridNotifier struct {
	EnvironmentVariable string
	DrilldownURL        string
	SendgridApiKey      string
	NoReplyName         string
	NoReplyEmail        string
	ContactName         string
	ContactEmail        string
}

func NewSendgridNotifier(config ServerConfig) Notifier {
	return &SendGridNotifier{
		EnvironmentVariable: config.NOTIFIER_ENVIRONMENT_VARIABLE,
		DrilldownURL:        config.HAPPENING_SERVER_URL,
		SendgridApiKey:      config.NOTIFIER_SENDGRID_API_KEY,
		NoReplyName:         config.NOTIFIER_NO_REPLY_NAME,
		NoReplyEmail:        config.NOTIFIER_NO_REPLY_EMAIL,
		ContactName:         config.NOTIFIER_CONTACT_NAME,
		ContactEmail:        config.NOTIFIER_CONTACT_EMAIL,
	}
}

func (notifier *SendGridNotifier) buildSendgridMail(notifierMail NotifierMail) []byte {
	from := mail.NewEmail(notifier.NoReplyName, notifier.NoReplyEmail)
	to := mail.NewEmail(notifier.ContactName, notifier.ContactEmail)
	replyTo := mail.NewEmail(notifier.ContactName, notifier.ContactEmail)
	subject := notifierMail.Subject()
	content := mail.NewContent("text/plain", notifierMail.Text())
	m := mail.NewV3MailInit(from, subject, to, content)
	m.SetReplyTo(replyTo)
	return mail.GetRequestBody(m)
}

func (notifier *SendGridNotifier) sendMail(notifierMail NotifierMail) {
	if notifier.SendgridApiKey == "" {
		log.Panicln("Sendgrid API key required in environment configuration")
	}
	request := sendgrid.GetRequest(
		notifier.SendgridApiKey,
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	request.Method = "POST"
	request.Body = notifier.buildSendgridMail(notifierMail)
	if _, err := sendgrid.API(request); err != nil {
		log.Panic(err)
	}
}

func (notifier *SendGridNotifier) Alert(check Check) {
	go notifier.sendMail(
		NotifierMail{
			Check:               check,
			EnvironmentVariable: notifier.EnvironmentVariable,
			DrilldownURL:        notifier.DrilldownURL,
		},
	)
}

func (notifier *SendGridNotifier) Resolve(check Check) {
	go notifier.sendMail(
		NotifierMail{
			Check:               check,
			EnvironmentVariable: notifier.EnvironmentVariable,
			DrilldownURL:        notifier.DrilldownURL,
			Resolved:            true,
		},
	)
}
