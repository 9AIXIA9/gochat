package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
	"gochat/internal/infra/snowflake"
	"gochat/internal/utils"
)

type Signup struct {
	repo domain.UserRepository
}

func NewSignup(repo domain.UserRepository) domain.SignupUsecase {
	return &Signup{repo: repo}
}

func (uc *Signup) Logic(ctx context.Context, req *domain.SignupRequest) (*domain.Message, error) {
	// 加密密码
	hashedPassword, err := uc.EncryptPwd(ctx, req.Password)
	if err != nil {
		return nil, err
	}

	// 生成用户号码
	userNumber, err := uc.GenerateNumber(ctx)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &domain.User{
		Number:  userNumber,
		Name:    req.Name,
		PwdHash: hashedPassword,
	}

	if err := uc.CreateUser(ctx, user); err != nil {
		if utils.IsDuplicate(err) {
			return domain.UserExistResponse, nil
		}
		return nil, err
	}

	// 返回用户信息
	return domain.NewSuccessMessage(domain.SignupResponse{UserNumber: userNumber}), nil
}

func (uc *Signup) EncryptPwd(ctx context.Context, pwd string) (string, error) {
	return encrypt.Encrypt(ctx, pwd)
}

func (uc *Signup) GenerateNumber(ctx context.Context) (domain.UserNumber, error) {
	return snowflake.GenerateUserNumber(ctx)
}

func (uc *Signup) CreateUser(ctx context.Context, user *domain.User) error {
	return uc.repo.Save(ctx, user)
}
