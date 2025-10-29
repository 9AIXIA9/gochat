package event

type Publisher interface {
	PublishEvents(events []Event) error
}
