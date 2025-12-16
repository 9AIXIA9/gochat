package application

import (
	"context"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
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
	accessTokenParser domain.AccessTokenParser
}

func NewParseAccessTokenUseCase(accessTokenParser domain.AccessTokenParser) (ParseAccessTokenUseCase, error) {
	if err := utils.CheckInterfaces(accessTokenParser); err != nil {
		return nil, err
	}

	return &parseAccessTokenUseCase{
		accessTokenParser: accessTokenParser,
	}, nil
}

func (uc *parseAccessTokenUseCase) Execute(_ context.Context, input *ParseAccessTokenInput) (*ParseAccessTokenOutput, error) {
	userID, err := uc.accessTokenParser.Parse(input.AccessToken)
	if err != nil {
		return nil, errors.WrapBusiness(err, "invalid access token")
	}
	return &ParseAccessTokenOutput{UserID: userID}, nil
}
