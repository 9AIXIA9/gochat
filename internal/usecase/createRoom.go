package usecase

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
	"gochat/internal/infra/snowflake"
	"strconv"
)

type CreateRoom struct {
	repo domain.RoomRepository
}

func NewCreateRoom(repo domain.RoomRepository) domain.CreateRoomUsecase {
	return &CreateRoom{repo: repo}
}

func (uc *CreateRoom) Logic(req *domain.CreateRoomRequest) (*domain.Response, error) {
	//加密secret
	secretHash, err := uc.EncryptSecret(req.Body.Secret)
	if err != nil {
		return nil, err
	}

	//生成房间号
	roomNumber := uc.GenerateNumber()

	// 创建房间
	room := &domain.Room{
		Name:         req.Body.Name,
		Number:       roomNumber,
		SecretHash:   secretHash,
		Description:  req.Body.Description,
		CurrentUsers: 1,
		MaxUsers:     req.Body.MaxUsers,
		Owner:        req.UserNumber,
	}

	if exist, err := uc.CreateRoom(room); err != nil {
		return nil, err
	} else if exist {
		return domain.NewResponseWithDefaultMsg(domain.CodeRoomExist), nil
	}

	// 返回房间信息
	return domain.NewSuccessResponse(gin.H{
		"room_number": strconv.Itoa(int(roomNumber)),
		"room_name":   req.Body.Name,
		"description": req.Body.Description,
		"max_users":   req.Body.MaxUsers,
		"owner":       strconv.Itoa(int(req.UserNumber)),
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
