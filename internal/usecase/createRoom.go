package usecase

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (uc *CreateRoom) Logic(req *domain.CreateRoomRequest) (*domain.Response, error) {
	//加密secret
	secretHash, err := uc.EncryptSecret(req.Secret)
	if err != nil {
		return nil, err
	}

	//生成房间号
	roomNumber := uc.GenerateNumber()

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

	if exist, err := uc.CreateRoom(room); err != nil {
		return nil, err
	} else if exist {
		return domain.NewResponseWithoutMsg(domain.CodeRoomExist), nil
	}

	zap.L().Info("create room successfully",
		zap.Int64("userNumber", int64(req.UserNumber)), zap.Int64("roomNumber", int64(roomNumber)))

	// 返回房间信息
	return domain.NewSuccessResponse(gin.H{
		"room_number": roomNumber,
		"room_name":   req.Name,
		"description": req.Description,
		"max_users":   req.MaxUsers,
		"owner":       req.UserNumber,
	}), nil
}

func (uc *CreateRoom) EncryptSecret(secret string) (string, error) {
	return encrypt.Encrypt(secret)
}

func (uc *CreateRoom) GenerateNumber() domain.RoomNumber {
	return snowflake.GenerateRoomNumber()
}

func (uc *CreateRoom) CreateRoom(room *domain.Room) (bool, error) {
	return uc.repo.Create(room)
}
