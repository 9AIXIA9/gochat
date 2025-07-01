package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
	"gochat/internal/infra/snowflake"
)

type CreateRoom struct {
	repo domain.RoomRepository
}

func NewCreateRoom(repo domain.RoomRepository) domain.CreateRoomUsecase {
	return &CreateRoom{repo: repo}
}

func (uc *CreateRoom) Logic(ctx context.Context, req *domain.CreateRoomRequest) (*domain.Message, error) {
	//加密secret
	secretHash, err := uc.EncryptSecret(ctx, req.Secret)
	if err != nil {
		return nil, err
	}

	//生成房间号
	roomNumber, err := uc.GenerateNumber(ctx)
	if err != nil {
		return nil, err
	}

	// 创建房间
	room := &domain.Room{
		Name:         req.Name,
		Number:       roomNumber,
		SecretHash:   secretHash,
		Description:  req.Description,
		CurrentUsers: 1,
		MaxUsers:     req.MaxUsers,
		Owner:        req.UserNumber,
	}

	if exist, err := uc.CreateRoom(ctx, room); err != nil {
		return nil, err
	} else if exist {
		return domain.RoomExistResponse, nil
	}

	// 返回房间信息
	return domain.NewSuccessMessage(domain.CreateRoomResponse{
		RoomNumber:  roomNumber,
		RoomName:    req.Name,
		Description: req.Description,
		MaxUsers:    req.MaxUsers,
		Owner:       req.UserNumber,
	}), nil
}

func (uc *CreateRoom) EncryptSecret(ctx context.Context, secret string) (string, error) {
	return encrypt.Encrypt(ctx, secret)
}

func (uc *CreateRoom) GenerateNumber(ctx context.Context) (domain.RoomNumber, error) {
	return snowflake.GenerateRoomNumber(ctx)
}

func (uc *CreateRoom) CreateRoom(ctx context.Context, room *domain.Room) (bool, error) {
	return uc.repo.Create(ctx, room)
}
