package http

import (
	"errors"
	"gochat/internal/chat/application"
	"gochat/internal/chat/dto"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const defaultRoomMessagesLimit = 20

type ListRoomMessagesRequest struct {
	UserID kernel.UserID    `json:"-" validate:"required"`
	BaseID kernel.MessageID `form:"base_id"`                                  // 用于分页游标
	Limit  int              `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
}

func (r *ListRoomMessagesRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultRoomMessagesLimit
	}
	return nil
}

type ListRoomMessagesResponseData struct {
	RoomMessages []*dto.RoomMessage `json:"room_messages,omitempty"`
}

func NewListRoomMessagesHandler(useCase application.ListRoomMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListRoomMessagesRequest) *application.ListRoomMessagesInput {
			return &application.ListRoomMessagesInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListRoomMessagesOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&ListRoomMessagesResponseData{
				RoomMessages: dto.ToRoomMessageDTOs(output.RoomMessages),
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			default:
				zap.L().Error("ListRoomMessageHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
