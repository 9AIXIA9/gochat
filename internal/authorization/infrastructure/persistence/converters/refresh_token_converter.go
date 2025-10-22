package converters

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/infrastructure/persistence/models"
	"gochat/internal/shared/kernel"
)

var _ kernel.GenericModelConverter[*models.RefreshToken, *domain.RefreshToken] = (*RefreshTokenConverter)(nil)

type RefreshTokenConverter struct {
}

func (c *RefreshTokenConverter) ToModel(refreshToken *domain.RefreshToken) (*models.RefreshToken, error) {
	return &models.RefreshToken{
		String:       refreshToken.String(),
		UserID:       refreshToken.UserID(),
		ExpiredAt:    refreshToken.ExpiredAt(),
		RefreshCount: refreshToken.RefreshCount(),
	}, nil
}

func (c *RefreshTokenConverter) ToDomain(refreshToken *models.RefreshToken) (*domain.RefreshToken, error) {
	return domain.NewRefreshToken(refreshToken.String, refreshToken.UserID, refreshToken.ExpiredAt, refreshToken.RefreshCount), nil
}
