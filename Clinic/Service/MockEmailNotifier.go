package service

// Notifier module
type Notifier[key comparable, value any] struct {
}

func (n *Notifier[k, v]) SendNotification(k, v) error {
	return nil
}

func ReturnNewEmailNotifier() *Notifier[string, string] {
	return &Notifier[string, string]{}
}
