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
	UserID kernel.UserID     `json:"-" validate:"required"`
	BaseID domain.RoomshipID `form:"base_id"`                                  // 用于分页游标
	Limit  int               `form:"limit" validate:"omitempty,min=1,max=100"` // 每页条数
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

// NewListRoomshipsHandler 获取房间关系列表
// @Summary      获取房间关系列表
// @Description  获取当前登录用户加入的房间关系列表，可基于 base_id 游标和 limit 分页
// @Tags         Roomship
// @Security     BearerAuth
// @Produce      json
// @Param        base_id  query       domain.RoomshipID     false "分页游标，返回该ID之前的记录"
// @Param        limit    query     int    false "分页大小，默认20，最大100"
// @Success      200      {object}  api.Response{data=ListRoomshipsResponseData} "成功返回房间关系列表"
// @Failure      400      {object}  api.Response   "请求参数错误"
// @Failure      401      {object}  api.Response   "未认证"
// @Failure      500      {object}  api.Response   "服务器内部错误"
// @Router       /roomship/ [get]
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
