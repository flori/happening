package happening

type NullNotifier struct {
}

func NewNullNotifier() Notifier {
	return &NullNotifier{}
}

func (notifier *NullNotifier) Alert(check Check) {
	return
}

func (notifier *NullNotifier) Resolve(check Check) {
	return
}
