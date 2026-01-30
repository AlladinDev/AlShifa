package interfaces

type INotifier[k comparable, v any] interface {
	SendNotification(key k, value k) error
}
