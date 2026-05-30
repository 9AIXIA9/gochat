//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package core

type SessionIDGenerator interface {
	Generate() SessionID
}
