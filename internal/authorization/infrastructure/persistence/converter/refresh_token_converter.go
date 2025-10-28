package converter

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
)

var _ gormutils.GenericModelConverter[*model.RefreshToken, *domain.RefreshTokenEntity] = (*RefreshTokenConverter)(nil)

type RefreshTokenConverter struct {
}

func (c *RefreshTokenConverter) ToModel(refreshToken *domain.RefreshTokenEntity) *model.RefreshToken {
	return &model.RefreshToken{
		Token:        refreshToken.Token(),
		UserID:       refreshToken.UserID(),
		ExpiredAt:    refreshToken.ExpiredAt(),
		RefreshCount: refreshToken.RefreshCount(),
	}
}

func (c *RefreshTokenConverter) ToDomain(refreshToken *model.RefreshToken) *domain.RefreshTokenEntity {
	return domain.NewRefreshToken(refreshToken.Token, refreshToken.UserID, refreshToken.ExpiredAt, refreshToken.RefreshCount)
}
