package http

import (
	"errors"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/dto"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const defaultFriendshipsLimit = 20

type ListFriendshipsRequest struct {
	UserID kernel.UserID       `json:"-" validate:"required"`
	BaseID domain.FriendshipID `form:"base_id"`                                  // 用于分页游标
	Limit  int                 `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
}

func (r *ListFriendshipsRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultFriendshipsLimit
	}
	return nil
}

type ListFriendshipsResponseData struct {
	Friendships []*dto.Friendship `json:"friendships,omitempty"`
}

func NewListFriendshipsHandler(useCase application.ListFriendshipsUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListFriendshipsRequest) *application.ListFriendshipsInput {
			return &application.ListFriendshipsInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListFriendshipsOutput) {
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&ListFriendshipsResponseData{
				Friendships: dto.ToFriendshipDTOs(output.Friendships),
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			default:
				zap.L().Error("ListFriendshipsHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
