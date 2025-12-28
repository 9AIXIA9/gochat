package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/dto"
	"gochat/internal/shared/kernel"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

const defaultRoomshipsLimit = 20

type ListRoomshipsRequest struct {
	UserID kernel.UserID     `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	BaseID domain.RoomshipID `form:"base_id"  validate:"omitempty" example:"019b593b-462e-74d6-bfda-0e103a172192"` // 用于分页游标
	Limit  int               `form:"limit" validate:"omitempty,min=1,max=100" example:"20"`                        // 每页条数
}

func (r *ListRoomshipsRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.ShouldBindQuery(r); err != nil {
		return err
	}
	// 默认值
	if r.Limit == 0 {
		r.Limit = defaultRoomshipsLimit
	}
	return nil
}

type ListRoomshipsResponseData struct {
	Roomships []*dto.Roomship `json:"roomships,omitempty"`
}

// NewListRoomshipsHandler 获取当前用户所在的房间关系列表
// @Summary      获取当前用户所在的房间关系列表
// @Description  获取当前登录用户加入的房间关系列表，可基于 base_id 游标和 limit 分页
// @Tags         Roomship
// @Security     BearerAuth
// @Param        base_id  query     string     false "分页游标，返回该ID之前的记录"     example("019b593b-462e-74d6-bfda-0e103a172190")
// @Param        limit    query     int    false "分页大小，默认20，最大100"   minimum(1) maximum(100) default(20) example(50)
// @Success      200      {object}  api.Response{data=ListRoomshipsResponseData} "成功返回房间关系列表"
// @Router       /rooms/ [get]
func NewListRoomshipsHandler(useCase application.ListRoomshipsUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *ListRoomshipsRequest) *application.ListRoomshipsInput {
			return &application.ListRoomshipsInput{
				UserID: request.UserID,
				BaseID: request.BaseID,
				Limit:  request.Limit,
			}
		},
		func(ginContext *gin.Context, output *application.ListRoomshipsOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &ListRoomshipsResponseData{
				Roomships: dto.ToRoomshipDTOs(output.Roomships),
			})
		},
	)
}
