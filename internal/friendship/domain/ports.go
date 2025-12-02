//go:generate mockgen -source=ports.go -destination=./mocks/mock_ports.go -package=mocks
package domain

import "gochat/internal/shared/kernel"

type FriendshipIDGenerator interface {
	Generate() FriendshipID
}

type OperationIDGenerator interface {
	Generate() kernel.OperationID
}
