package happening

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"strings"
)

type SMTPNotifier struct {
	EnvironmentVariable string
	DrilldownURL        string
	SMTPServerURL       string
	NoReplyName         string
	NoReplyEmail        string
	ContactName         string
	ContactEmail        string
	SkipSSLVerify       bool
}

func NewSMTPNotifier(config ServerConfig) Notifier {
	return &SMTPNotifier{
		EnvironmentVariable: config.NOTIFIER_ENVIRONMENT_VARIABLE,
		DrilldownURL:        config.HAPPENING_SERVER_URL,
		SMTPServerURL:       config.NOTIFIER_SMTP_SERVER_URL,
		NoReplyName:         config.NOTIFIER_NO_REPLY_NAME,
		NoReplyEmail:        config.NOTIFIER_NO_REPLY_EMAIL,
		ContactName:         config.NOTIFIER_CONTACT_NAME,
		ContactEmail:        config.NOTIFIER_CONTACT_EMAIL,
	}
}

func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func performSendMail(addr, serverName string, a smtp.Auth, from string, to string, msg []byte, skipVerify bool) error {
	if err := validateLine(from); err != nil {
		return err
	}
	if err := validateLine(to); err != nil {
		return err
	}
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{
			ServerName:         serverName,
			InsecureSkipVerify: skipVerify,
		}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func (notifier *SMTPNotifier) sendMail(notifierMail NotifierMail) {
	serverURL, err := url.Parse(notifier.SMTPServerURL)
	if err != nil {
		log.Panic(err)
	}

	skipVerify := false
	if serverURL.Query().Get("sslVerify") == "disable" {
		skipVerify = true
	}

	password, _ := serverURL.User.Password()
	auth := smtp.PlainAuth("", serverURL.User.Username(), password, serverURL.Hostname())

	toAddr := fmt.Sprintf("%s <%s>", notifier.ContactName, notifier.ContactEmail)
	fromAddr := fmt.Sprintf("%s <%s>", notifier.NoReplyName, notifier.NoReplyEmail)
	mail := fmt.Sprintf("From: %s\r\n", fromAddr)
	mail += fmt.Sprintf("To: %s\r\n", toAddr)
	mail += fmt.Sprintf("Subject: %s\r\n", notifierMail.Subject())
	mail += "\r\n"
	mail += notifierMail.Text()
	mail += "\r\n"

	err = performSendMail(
		serverURL.Host,
		serverURL.Hostname(),
		auth,
		notifier.NoReplyEmail,
		notifier.ContactEmail,
		[]byte(mail),
		skipVerify,
	)
	if err != nil {
		log.Printf("Sending mail to %s failed: %v", serverURL.Host, err)
	}
}

func (notifier *SMTPNotifier) Alert(check Check) {
	go notifier.sendMail(
		NotifierMail{
			Check:               check,
			EnvironmentVariable: notifier.EnvironmentVariable,
			DrilldownURL:        notifier.DrilldownURL,
		},
	)
}

func (notifier *SMTPNotifier) Resolve(check Check) {
	go notifier.sendMail(
		NotifierMail{
			Check:               check,
			EnvironmentVariable: notifier.EnvironmentVariable,
			DrilldownURL:        notifier.DrilldownURL,
			Resolved:            true,
		},
	)
}
