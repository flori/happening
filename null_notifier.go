package happening

type NullNotifier struct {
}

func NewNullNotifier() Notifier {
	return &NullNotifier{}
}

func (notifier *NullNotifier) Alert(check Check) {
	return
}
