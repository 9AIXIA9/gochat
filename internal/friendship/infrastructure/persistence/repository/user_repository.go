package repository

import (
	"context"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
		DoNothing: true,
	}).Create(&model.User{
		ID:     user.ID(),
		Number: user.Number(),
	}).Error)
}

func (repo *UserRepository) FindByNumber(ctx context.Context, number kernel.UserNumber) (*domain.User, error) {
	var user model.User
	err := repo.db.WithContext(ctx).First(&user, "number = ?", number).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}

	//TODO 暂时不加载好友和请求
	return domain.LoadUser(user.ID, user.Number, make([]kernel.UserID, 0), make([]*domain.FriendRequest, 0)), nil
}
