package happening

type NotifierMail interface {
	Subject() string
	Text() string
}
