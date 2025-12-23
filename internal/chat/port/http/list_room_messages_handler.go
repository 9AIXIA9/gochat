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
	UserID kernel.UserID    `json:"-" validate:"required"`
	RoomID kernel.RoomID    `uri:"room_id" validate:"required"`
	BaseID kernel.MessageID `form:"base_id"`                                  // 用于分页游标
	Limit  int              `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
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
// @Produce      json
// @Param        room_id  path      kernel.RoomID       true  "房间ID"
// @Param        base_id  query     kernel.MessageID     false "分页游标，返回该ID之前的消息"
// @Param        limit    query     int    false "分页大小，默认20，最大100"
// @Success      200      {object}  api.Response{data=ListRoomMessagesResponseData} "成功返回房间消息记录"
// @Failure      400      {object}  api.Response "请求参数错误"
// @Failure      401      {object}  api.Response "未认证"
// @Failure      500      {object}  api.Response "服务器内部错误"
// @Router       /chats/rooms/messages [get]
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
