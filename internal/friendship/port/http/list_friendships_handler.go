package http

import (
	"errors"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/dto"
	ginutils "gochat/internal/infrastructure/gin"
	sharedHttp "gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

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

// NewListFriendshipsHandler 获取好友列表
// @Summary      获取好友列表
// @Description  获取当前登录用户的好友关系列表，可基于 base_id 游标和 limit 分页
// @Tags         Friendship
// @Security     BearerAuth
// @Produce      json
// @Param        base_id  query     int    false "分页游标，返回该ID之前的记录"
// @Param        limit    query     int    false "分页大小，默认20，最大100"
// @Success      200      {object}  ListFriendshipsResponseData "成功返回好友列表"
// @Failure      400      {object}  sharedHttp.Response      "请求参数错误"
// @Failure      401      {object}  sharedHttp.Response      "未认证"
// @Failure      500      {object}  sharedHttp.Response      "服务器内部错误"
// @Router       /friendship/ [get]
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
			ginutils.ResponseSuccessWithData(ginContext, &ListFriendshipsResponseData{
				Friendships: dto.ToFriendshipDTOs(output.Friendships),
			})
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			default:
				zap.L().Error("ListFriendshipsHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
	)
}
