package uuid

import (
	"gochat/internal/friendship/domain"

	"github.com/google/uuid"
)

var _ domain.FriendshipIDGenerator = (*FriendshipIDGenerator)(nil)

type FriendshipIDGenerator struct{}

func NewFriendshipIDGenerator() *FriendshipIDGenerator {
	return &FriendshipIDGenerator{}
}

func (FriendshipIDGenerator) Generate() domain.FriendshipID {
	return domain.FriendshipID(uuid.Must(uuid.NewV7()).String())
}
