//go:generate mockgen -source=publisher.go -destination=./mocks/mock_publisher.go -package=mocks
package event

type Publisher interface {
	Publish(event Event) error
}
