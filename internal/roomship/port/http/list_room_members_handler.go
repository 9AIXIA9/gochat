package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/dto"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListRoomMembersRequest struct {
	RoomID kernel.RoomID `uri:"room_id" validate:"required"`
}

func (r *ListRoomMembersRequest) Bind(ginContext *gin.Context) error {
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

type ListRoomMembersResponseData struct {
	Roomships []*dto.Roomship `json:"roomships,omitempty"`
}

// NewListRoomMembersHandler 获取房间成员列表
// @Summary      获取房间成员列表
// @Description  根据房间ID获取房间成员列表
// @Tags         Roomship
// @Produce      json
// @Param        room_id  path        kernel.RoomID                           true  "房间ID"
// @Success      200      {object}  ListRoomMembersResponseData  "成功返回房间成员列表"
// @Failure      400      {object}  api.Response       "请求参数错误"
// @Failure      404      {object}  api.Response       "房间不存在"
// @Failure      500      {object}  api.Response       "服务器内部错误"
// @Router       /roomship/room/{room_id} [get]
func NewListRoomMembersHandler(useCase application.ListRoomMembersUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListRoomMembersRequest) *application.ListRoomMembersInput {
			return &application.ListRoomMembersInput{
				RoomID: request.RoomID,
			}
		},
		func(ginContext *gin.Context, output *application.ListRoomMembersOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &ListRoomMembersResponseData{
				Roomships: dto.ToRoomshipDTOs(output.Roomships),
			})
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "input is empty")
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.ResponseWithMessage(ginContext, api.CodeNotFound, "room is not found")
			default:
				zap.L().Error("ListRoomMembersHandler error", zap.Error(err))
				ginutils.Response(ginContext, api.CodeServerError)
			}
		},
	)
}
