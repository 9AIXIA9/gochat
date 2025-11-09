package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	myErrors "gochat/internal/shared/errors"

	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
	"gochat/internal/social/domain"
	"gochat/internal/social/infrastructure/persistence/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ application.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Save upsert 房间并同步成员关系
func (repo *RoomRepository) Save(ctx context.Context, room *domain.Room) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Upsert room basic fields (不含 Members)
		if err := tx.Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"owner", "number", "password_encrypted", "member_count", "max_member_count"}),
			},
		).Create(toModelRoom(room)).Error; err != nil {
			return err
		}

		// 确保取到持久化后的房间实体用于关联操作
		var persisted model.Room
		if err := tx.First(&persisted, "id = ?", room.ID()).Error; err != nil {
			return err
		}

		// 同步成员关系到多对多表
		memberUsers, err := repo.findUsersByIDsTx(tx, room.Members())
		if err != nil {
			return err
		}
		// Replace 将以传入集合为准，原子同步
		if err := tx.Model(&persisted).Association("Members").Replace(memberUsers); err != nil {
			return err
		}

		return nil
	})
}

// FindByNumber 通过房间号查询
func (repo *RoomRepository) FindByNumber(ctx context.Context, number domain.RoomNumber) (*domain.Room, error) {
	var m model.Room
	err := repo.db.WithContext(ctx).
		Preload("Members").
		Where("number = ?", number.String()).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return toDomainRoom(&m), nil
}

// Join 将用户加入房间（原子更新成员和计数）
func (repo *RoomRepository) Join(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var r model.Room
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&r, "id = ?", roomID.String()).Error; err != nil {
			return err
		}

		// 检查是否已加入
		var cnt int64
		if err := tx.Table("gochat.social_room_members").
			Where("room_id = ? AND user_id = ?", roomID.String(), userID).
			Count(&cnt).Error; err != nil {
			return err
		}
		if cnt > 0 {
			return nil
		}

		// 检查人数上限
		if r.MaxMemberCount > 0 && r.MemberCount >= r.MaxMemberCount {
			return myErrors.ErrExceedMaxValue
		}

		// 确认用户存在
		var u model.User
		if err := tx.First(&u, "id = ?", userID).Error; err != nil {
			return err
		}

		// 建立关联
		if err := tx.Model(&r).Association("Members").Append(&u); err != nil {
			return err
		}

		// 增加计数
		if err := tx.Model(&r).UpdateColumn("member_count", gorm.Expr("member_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	}))
}

// Leave 将用户从房间移除（原子更新成员和计数）
func (repo *RoomRepository) Leave(ctx context.Context, roomID domain.RoomID, userID kernel.UserID) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var r model.Room
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&r, "id = ?", roomID.String()).Error; err != nil {
			return err
		}

		// 拥有者不能退出
		if r.Owner == userID {
			return myErrors.ErrOwnerCantLeave
		}

		// 检查是否成员
		var cnt int64
		if err := tx.Table("gochat.social_room_members").
			Where("room_id = ? AND user_id = ?", roomID.String(), userID).
			Count(&cnt).Error; err != nil {
			return err
		}
		if cnt == 0 {
			// 非成员，幂等返回
			return nil
		}

		// 删除关联
		if err := tx.Model(&r).Association("Members").Delete(&model.User{ID: userID}); err != nil {
			return err
		}

		// 减少计数（不小于 0）
		if err := tx.Model(&r).Where("member_count > 0").
			UpdateColumn("member_count", gorm.Expr("member_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	}))
}

// FindByID 通过房间ID查询
func (repo *RoomRepository) FindByID(ctx context.Context, id domain.RoomID) (*domain.Room, error) {
	var m model.Room
	err := repo.db.WithContext(ctx).
		Preload("Members").
		First(&m, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}
	return toDomainRoom(&m), nil
}

func (repo *RoomRepository) findUsersByIDsTx(tx *gorm.DB, ids []kernel.UserID) ([]*model.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var users []*model.User
	if err := tx.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	// 可选：严格校验全部存在
	if len(users) != len(ids) {
		return nil, myErrors.ErrNotFound
	}
	return users, nil
}

func toModelRoom(room *domain.Room) *model.Room {
	return &model.Room{
		ID:                room.ID(),
		Owner:             room.Owner(),
		Number:            room.Number(),
		PasswordEncrypted: room.PasswordEncrypted(),
		MemberCount:       room.MemberCount(),
		MaxMemberCount:    room.MaxMemberCount(),
	}
}

func toDomainRoom(m *model.Room) *domain.Room {
	memberIDs := make([]kernel.UserID, 0, len(m.Members))
	for _, u := range m.Members {
		memberIDs = append(memberIDs, u.ID)
	}
	return domain.NewRoom(
		m.ID,
		m.Owner,
		m.Number,
		m.PasswordEncrypted,
		memberIDs,
		m.MemberCount,
		m.MaxMemberCount,
		m.CreatedAt,
	)
}
