package presentation

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gochat/internal/domain"
	"reflect"
	"strings"
)

func BindParams(c *gin.Context, param interface{}) (*domain.Response, error) {
	// 校验空指针
	if param == nil || reflect.ValueOf(param).Kind() != reflect.Ptr || reflect.ValueOf(param).IsNil() {
		return nil, errors.New("param is a nil pointer")
	}

	// 获取并校验参数
	if err := c.ShouldBind(param); err != nil {
		var typeErr validator.ValidationErrors
		if errors.As(err, &typeErr) {
			// 使用 strings.Builder 构建错误消息
			msg := strings.Builder{}
			for i, fe := range typeErr {
				if i > 0 {
					msg.WriteString("; ")
				}
				fieldName := fe.Field()
				trans := GetTrans()
				t, terr := trans.T(fe.Tag(), fieldName, fe.Param())
				if terr != nil {
					msg.WriteString(fieldName + ": " + fe.Error())
				} else {
					msg.WriteString(fieldName + ": " + t)
				}
			}
			return domain.NewResponse(domain.CodeInvalidParam, msg.String()), nil
		} else {
			return domain.NewResponse(domain.CodeInvalidParam, err.Error()), nil
		}
	}
	return nil, nil
}
