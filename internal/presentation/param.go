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

	// 获取参数类型
	t := reflect.TypeOf(param).Elem()

	// 绑定URI参数(如果存在URI字段)
	if hasField(t, "URI") {
		uriValue := reflect.ValueOf(param).Elem().FieldByName("URI").Addr().Interface()
		if err := c.ShouldBindUri(uriValue); err != nil {
			var validErr validator.ValidationErrors
			if errors.As(err, &validErr) {
				msg := buildValidationErrorMessage(validErr)
				return domain.NewResponse(domain.CodeInvalidParam, msg), nil
			}
			return nil, err
		}
	}

	// 绑定Body参数(如果存在Body字段)
	if hasField(t, "Body") {
		bodyValue := reflect.ValueOf(param).Elem().FieldByName("Body").Addr().Interface()
		if err := c.ShouldBindJSON(bodyValue); err != nil {
			var validErr validator.ValidationErrors
			if errors.As(err, &validErr) {
				msg := buildValidationErrorMessage(validErr)
				return domain.NewResponse(domain.CodeInvalidParam, msg), nil
			}
			return nil, err
		}
	}

	return nil, nil
}

// hasField 检查结构体是否包含特定字段
func hasField(t reflect.Type, fieldName string) bool {
	if t.Kind() != reflect.Struct {
		return false
	}

	_, found := t.FieldByName(fieldName)
	return found
}

// 提取验证错误信息格式化为辅助函数
func buildValidationErrorMessage(typeErr validator.ValidationErrors) string {
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
	return msg.String()
}
