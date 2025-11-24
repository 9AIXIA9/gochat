package application

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/kernel"
)

type ParseAccessTokenUseCase kernel.UseCase[*ParseAccessTokenInput, *ParseAccessTokenOutput]

type ParseAccessTokenInput struct {
	AccessToken domain.AccessToken
}

func (input *ParseAccessTokenInput) Validate() error {
	return input.AccessToken.Validate()
}

type ParseAccessTokenOutput struct {
	UserID kernel.UserID
}

type parseAccessTokenUseCase struct {
	accessTokenParser AccessTokenParser
}

func NewParseAccessTokenUseCase(accessTokenParser AccessTokenParser) ParseAccessTokenUseCase {
	return &parseAccessTokenUseCase{accessTokenParser: accessTokenParser}
}

func (uc *parseAccessTokenUseCase) Execute(_ context.Context, input *ParseAccessTokenInput) (*ParseAccessTokenOutput, error) {
	userID, err := uc.accessTokenParser.Parse(input.AccessToken)
	if err != nil {
		return nil, err
	}
	return &ParseAccessTokenOutput{UserID: userID}, nil
}
