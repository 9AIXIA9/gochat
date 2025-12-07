package model

import (
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/kernel"
	"time"
)

//TODO 为所有模型添加索引，外键约束

type FriendRequest struct {
	ID      kernel.OperationID        `gorm:"primaryKey;type:char(36)"`
	From    kernel.UserID             `gorm:"not null;index:idx_friend_req_from_to_state,priority:1;index:idx_friend_req_from;"`
	To      kernel.UserID             `gorm:"not null;index:idx_friend_req_from_to_state,priority:2;index:idx_friend_req_to;"`
	Content string                    `gorm:"type:varchar(500)"`
	State   domain.FriendRequestState `gorm:"not null;index:idx_friend_req_from_to_state,priority:3;index:idx_friend_req_state"`
	SentAt  time.Time                 `gorm:"not null;index:idx_friend_req_to_sentat,priority:2;index:idx_friend_req_sentat"`
}

func (*FriendRequest) TableName() string {
	return "friendship_friend_requests"
}
