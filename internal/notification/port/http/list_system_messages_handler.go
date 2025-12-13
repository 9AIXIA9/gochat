package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/notification/application"
	"gochat/internal/notification/dto"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const defaultSystemMessagesLimit = 20

type ListSystemMessagesRequest struct {
	UserID kernel.UserID    `json:"-" validate:"required"`
	BaseID kernel.MessageID `form:"base_id"`                                  // 用于分页游标
	Limit  int              `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
}

func (r *ListSystemMessagesRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultSystemMessagesLimit
	}
	return nil
}

type ListSystemMessagesResponseData struct {
	SystemMessages []*dto.SystemMessage `json:"system_messages,omitempty"`
}

func NewListSystemMessagesHandler(useCase application.ListSystemMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListSystemMessagesRequest) *application.ListSystemMessagesInput {
			return &application.ListSystemMessagesInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListSystemMessagesOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &ListSystemMessagesResponseData{
				SystemMessages: dto.ToSystemMessageDTOs(output.SystemMessages),
			})
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			default:
				zap.L().Error("ListSystemMessagesHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
