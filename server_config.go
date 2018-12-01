package happening

type ServerConfig struct {
	PORT                          string `default:"8080"`
	DATABASE_NAME                 string `default:"happening"`
	POSTGRES_URL                  string `default:"postgresql://flori@localhost:5432/%s?sslmode=disable"`
	HTTP_REALM                    string `default:"happening"`
	HTTP_AUTH                     string
	NOTIFIER_KIND                 string `default:"MailCommand"`
	NOTIFIER_ENVIRONMENT_VARIABLE string `default:"RAILS_ENV"`
	NOTIFIER_NO_REPLY_NAME        string `default:"Happening"`
	NOTIFIER_NO_REPLY_EMAIL       string `default:"no-reply@localhost"`
	NOTIFIER_CONTACT_NAME         string `default:"Root"`
	NOTIFIER_CONTACT_EMAIL        string `default:"root@localhost"`
	NOTIFIER_SENDGRID_API_KEY     string
	NOTIFIER_MAIL_COMMAND         string `default:"mail"`
}
