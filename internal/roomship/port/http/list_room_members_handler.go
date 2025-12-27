package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/dto"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
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
// @Param        room_id  path        kernel.RoomID                           true  "房间ID"
// @Success      200      {object}  api.Response{data=ListRoomMembersResponseData}  "成功返回房间成员列表"
// @Router       /rooms/{room_id}/members [get]
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
	)
}
