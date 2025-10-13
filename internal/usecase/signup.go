package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/utils"
)

type Signup struct {
	domain.UserSaver
	domain.NumberGenerator
	domain.Encryptor
}

func NewSignup(saver domain.UserSaver,
	generator domain.NumberGenerator,
	encryptor domain.Encryptor) domain.SignupUsecase {
	return &Signup{
		UserSaver:       saver,
		NumberGenerator: generator,
		Encryptor:       encryptor,
	}
}

func (uc *Signup) Execute(ctx context.Context, req *domain.SignupRequest) (*domain.Response, error) {
	// 加密密码
	hashedPassword, err := uc.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	// 生成用户号码
	userNumber := uc.GenerateNumber()

	// 创建用户
	user := domain.CreateUser(userNumber, hashedPassword)

	if err := uc.SaveUser(ctx, user); err != nil {
		if utils.IsDuplicate(err) {
			return domain.UserExistResponse, nil
		}
		return nil, err
	}

	// 返回用户信息
	return domain.NewSuccessResponse(domain.SignupResponse{UserNumber: user.Number()}), nil
}
