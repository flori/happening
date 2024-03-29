package happening

type ServerConfig struct {
	PORT                          string `default:"8080"`
	POSTGRES_URL                  string `default:"postgresql://flori@localhost:5432/happening?sslmode=disable"`
	HAPPENING_SERVER_URL          string `default:"http://localhost:8080"`
	HTTP_AUTH                     string
	SIGNING_SECRET                string
	WEBUI_DIR                     string `default:"/webui"`
	NOTIFIER_KIND                 string
	NOTIFIER_ENVIRONMENT_VARIABLE string `default:"RAILS_ENV"`
	NOTIFIER_NO_REPLY_NAME        string `default:"Happening"`
	NOTIFIER_NO_REPLY_EMAIL       string `default:"no-reply@localhost"`
	NOTIFIER_CONTACT_NAME         string `default:"Root"`
	NOTIFIER_CONTACT_EMAIL        string `default:"root@localhost"`
	NOTIFIER_SENDGRID_API_KEY     string
	NOTIFIER_MAIL_COMMAND         string `default:"mail"`
	NOTIFIER_SMTP_SERVER_URL      string `default:"smtp://username:password@localhost:25?sslVerify=disable"`
	RETENTION_DAYS                string `default:"180"`
}
