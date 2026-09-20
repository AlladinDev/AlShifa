package interfaces

type INotifier interface {
	SendNotification(channel string, message string, title string, info string) error
}
