package converter

import (
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
	"gochat/internal/social/infrastructure/persistence/model"

	"gorm.io/gorm"
)

var _ gormutils.GenericModelConverter[*model.Room, *domain.Room] = (*RoomConverter)(nil)

type RoomConverter struct {
}

func (c *RoomConverter) ToModel(room *domain.Room) *model.Room {
	domainMembers := room.Members()
	members := make([]*model.RoomMember, 0, len(domainMembers))
	for _, member := range domainMembers {
		members = append(members, &model.RoomMember{
			RoomID: room.ID(),
			Member: member,
		})
	}

	return &model.Room{
		Model:             gorm.Model{},
		ID:                room.ID(),
		Owner:             room.Owner(),
		Number:            room.Number(),
		PasswordEncrypted: room.PasswordEncrypted(),
		MemberCount:       room.MemberCount(),
		MaxMemberCount:    room.MaxMemberCount(),
		CreatedAt:         room.CreatedAt(),
		Members:           members,
	}
}

func (c *RoomConverter) ToDomain(room *model.Room) *domain.Room {
	members := make([]kernel.UserID, 0, len(room.Members))
	for _, member := range room.Members {
		members = append(members, member.Member)
	}
	return domain.NewRoom(
		room.ID,
		room.Owner,
		room.Number,
		room.PasswordEncrypted,
		members,
		room.MemberCount,
		room.MaxMemberCount,
		room.CreatedAt,
	)
}
