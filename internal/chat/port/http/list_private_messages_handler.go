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

const defaultPrivateMessagesLimit = 20

type ListPrivateMessagesRequest struct {
	UserID kernel.UserID    `json:"-" validate:"required"`
	BaseID kernel.MessageID `form:"base_id"`                                  // 用于分页游标
	Limit  int              `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
}

func (r *ListPrivateMessagesRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultPrivateMessagesLimit
	}
	return nil
}

type ListPrivateMessagesResponseData struct {
	PrivateMessages []*dto.PrivateMessage `json:"private_messages,omitempty"`
}

func NewListPrivateMessagesHandler(useCase application.ListPrivateMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListPrivateMessagesRequest) *application.ListPrivateMessagesInput {
			return &application.ListPrivateMessagesInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListPrivateMessagesOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&ListPrivateMessagesResponseData{
				PrivateMessages: dto.ToPrivateMessageDTOs(output.PrivateMessages),
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			default:
				zap.L().Error("SendPrivateMessageHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
