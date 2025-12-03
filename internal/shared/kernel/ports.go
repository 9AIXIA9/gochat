//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package kernel

type MessageIDGenerator interface {
	Generate() MessageID
}
