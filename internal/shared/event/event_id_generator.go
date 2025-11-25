//go:generate mockgen -source=event_id_generator.go -destination=./mocks/mock_event_id_generator.go -package=mocks
package event

type IDGenerator interface {
	Generate() ID
}
