package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/dto"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

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
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&ListRoomMembersResponseData{
				Roomships: dto.ToRoomshipDTOs(output.Roomships),
			}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "input is empty"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeNotFound, "room is not found"))
			default:
				zap.L().Error("ListRoomMembersHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
