package useCase

import (
	"context"
	"gochat/internal/shared/kernel"
)

type ParseAccessToken struct {
	accessTokenParser kernel.AccessTokenParser
}

func NewParseAccessToken(accessTokenParser kernel.AccessTokenParser) kernel.ParseAccessTokenUseCase {
	return &ParseAccessToken{accessTokenParser: accessTokenParser}
}

func (uc *ParseAccessToken) Execute(ctx context.Context, input *kernel.ParseAccessTokenInput) (*kernel.ParseAccessTokenOutput, error) {
	userID, err := uc.accessTokenParser.Parse(input.AccessToken)
	if err != nil {
		return nil, err
	}
	return &kernel.ParseAccessTokenOutput{UserID: userID}, nil
}
