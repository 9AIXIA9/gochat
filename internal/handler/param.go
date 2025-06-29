package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gochat/internal/domain"
	"reflect"
	"strings"
)

func BindParams(c *gin.Context, param interface{}) (*domain.Response, error) {
	// 校验参数是否有效
	if param == nil {
		return nil, errors.New("param is nil")
	}

	// 绑定URI参数
	if err := bindUriIfExists(c, param); err != nil {
		return handleBindError(err)
	}

	// 绑定Body参数
	if err := bindBodyIfExists(c, param); err != nil {
		return handleBindError(err)
	}

	return nil, nil
}

func bindUriIfExists(c *gin.Context, param interface{}) error {
	if hasStructField(param, "URI") {
		return c.ShouldBindUri(getFieldAddr(param, "URI"))
	}
	return nil
}

func bindBodyIfExists(c *gin.Context, param interface{}) error {
	if hasStructField(param, "Body") {
		return c.ShouldBindJSON(getFieldAddr(param, "Body"))
	}
	return nil
}

// getFieldAddr 获取结构体中指定字段的地址
func getFieldAddr(obj interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(obj)

	// 确保我们处理的是指针类型
	if v.Kind() != reflect.Ptr {
		return nil
	}

	// 获取指针指向的结构体
	v = v.Elem()

	// 确保是结构体
	if v.Kind() != reflect.Struct {
		return nil
	}

	// 获取字段值
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	// 返回字段地址
	return field.Addr().Interface()
}

// hasStructField 检查结构体是否包含特定字段
func hasStructField(obj interface{}, fieldName string) bool {
	t := reflect.TypeOf(obj)

	// 确保我们处理的是指针类型
	if t.Kind() != reflect.Ptr {
		return false
	}

	// 获取指针指向的结构体的类型
	t = t.Elem()

	// 确保是结构体
	if t.Kind() != reflect.Struct {
		return false
	}

	// 查找字段
	_, found := t.FieldByName(fieldName)
	return found
}

func handleBindError(err error) (*domain.Response, error) {
	var validErr validator.ValidationErrors
	if errors.As(err, &validErr) {
		msg := buildValidationErrorMessage(validErr)
		return domain.NewResponse(domain.CodeInvalidParam, msg), nil
	}
	return nil, err
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
