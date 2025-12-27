package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type CreateRoomRequest struct {
	UserID         kernel.UserID   `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	MaxMemberCount int             `json:"max_member_count" validate:"gte=0,lte=200" example:"20"`
	Password       domain.Password `json:"password" validate:"max=100" example:"secret"`
}

func (r *CreateRoomRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.BindJSON(r); err != nil {
		return err
	}
	return nil
}

// NewCreateRoomHandler 创建房间
// @Summary      创建房间
// @Description  创建一个新的房间，当前登录用户将作为房主加入房间
// @Tags         Roomship
// @Security     BearerAuth
// @Param        request  body      CreateRoomRequest       true  "创建房间请求体"
// @Success      200      {object}  api.Response  "创建成功"
// @Router       /rooms/ [post]
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
	)
}
