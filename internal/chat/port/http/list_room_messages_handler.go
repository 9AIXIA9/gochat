package http

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/dto"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

const defaultRoomMessagesLimit = 20

type ListRoomMessagesRequest struct {
	UserID kernel.UserID    `json:"-" validate:"required" example:"019b5929-65bc-7549-89c6-3f7dc6872577"`
	RoomID kernel.RoomID    `uri:"room_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172190"`
	BaseID kernel.MessageID `form:"base_id" validate:"omitempty" example:"019b593b-462e-74d6-bfda-0e103a172190"` // 用于分页游标
	Limit  int              `form:"limit" validate:"omitempty,min=1,max=100" example:"20"`                       // 每页条数
}

func (r *ListRoomMessagesRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	if err := ginContext.BindUri(r); err != nil {
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

// NewListRoomMessagesHandler 获取房间内的消息
// @Summary      获取房间内的消息
// @Description  获取当前登录用户所在房间的消息记录，可基于 base_id 游标和 limit 分页
// @Tags         Chat
// @Security     BearerAuth
// @Param        room_id  path      string  true  "群聊对象房间ID"  example(019b5929-65bc-7549-89c6-3f7dc6872577)
// @Param        base_id  query     string  false "分页游标"       example(019b593b-462e-74d6-bfda-0e103a172190)
// @Param        limit    query     int     false "分页大小，默认 20，最大 100"      minimum(1) maximum(100) default(20) example(50)
// @Success      200      {object}  api.Response{data=ListRoomMessagesResponseData} "成功返回房间消息记录"
// @Router       /chats/rooms/messages/{room_id} [get]
func NewListRoomMessagesHandler(useCase application.ListRoomMessagesUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListRoomMessagesRequest) *application.ListRoomMessagesInput {
			return &application.ListRoomMessagesInput{
				UserID: request.UserID,
				RoomID: request.RoomID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListRoomMessagesOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &ListRoomMessagesResponseData{
				RoomMessages: dto.ToRoomMessageDTOs(output.RoomMessages),
			})
		},
	)
}
