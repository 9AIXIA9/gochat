package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateRoomRequest struct {
	UserID         kernel.UserID   `json:"-" validate:"required"`
	MaxMemberCount int             `json:"max_member_count" validate:"gte=0,lte=200"`
	Password       domain.Password `json:"password" validate:"max=100"`
}

func (r *CreateRoomRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBind(r); err != nil {
		return err
	}
	return nil
}

// NewCreateRoomHandler 创建房间
// @Summary      创建房间
// @Description  创建一个新的房间，当前登录用户将作为房主加入房间
// @Tags         Roomship
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      CreateRoomRequest       true  "创建房间请求体"
// @Success      201      {object}  sharedHttp.ApiResponse  "创建成功"
// @Failure      400      {object}  sharedHttp.ApiResponse  "请求参数错误"
// @Failure      401      {object}  sharedHttp.ApiResponse  "未认证"
// @Failure      500      {object}  sharedHttp.ApiResponse  "服务器内部错误"
// @Router       /roomship/room [post]
func NewCreateRoomHandler(useCase application.CreateRoomUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *CreateRoomRequest) *application.CreateRoomInput {
			return &application.CreateRoomInput{
				UserID:         request.UserID,
				MaxMemberCount: request.MaxMemberCount,
				Password:       request.Password,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			default:
				zap.L().Error("CreateRoomHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
