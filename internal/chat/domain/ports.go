//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

import "gochat/internal/shared/kernel"

type MessageIDGenerator interface {
	Generate() kernel.MessageID
}
