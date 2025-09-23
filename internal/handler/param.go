package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gochat/internal/domain"
	"gochat/internal/types"
	"gochat/internal/utils"
	"reflect"
	"strings"
)

func BindParams(c *gin.Context, param interface{}) (*domain.Response, error) {
	// 校验参数是否有效
	if param == nil {
		return nil, types.ErrNullPointer
	}

	// 绑定URI参数
	if err := bindUriIfExists(c, param); err != nil {
		return handleBindError(err)
	}

	// 绑定Query参数
	if err := bindQueryIfExists(c, param); err != nil {
		return handleBindError(err)
	}

	// 绑定Body参数
	if err := bindBodyIfExists(c, param); err != nil {
		return handleBindError(err)
	}

	// 绑定Cookie参数
	if err := bindCookieIfExists(c, param); err != nil {
		return handleBindError(err)
	}

	return nil, nil
}

func bindUriIfExists(c *gin.Context, param interface{}) error {
	if utils.HasStructField(param, "URI") {
		return c.ShouldBindUri(utils.GetFieldAddr(param, "URI"))
	}
	return nil
}

func bindQueryIfExists(c *gin.Context, param interface{}) error {
	if utils.HasStructField(param, "Query") {
		return c.ShouldBindQuery(utils.GetFieldAddr(param, "Query"))
	}
	return nil
}

func bindBodyIfExists(c *gin.Context, param interface{}) error {
	if utils.HasStructField(param, "Body") {
		return c.ShouldBindJSON(utils.GetFieldAddr(param, "Body"))
	}
	return nil
}

func bindCookieIfExists(c *gin.Context, param interface{}) error {
	if !utils.HasStructField(param, "Cookie") {
		return nil
	}
	cookiePtr := utils.GetFieldAddr(param, "Cookie")
	if cookiePtr == nil {
		return nil
	}
	cookieVal := reflect.ValueOf(cookiePtr).Elem()
	cookieType := cookieVal.Type()
	for i := 0; i < cookieType.NumField(); i++ {
		field := cookieType.Field(i)
		cookieName := field.Tag.Get("json")
		if cookieName == "" {
			cookieName = field.Name
		}
		val, err := c.Cookie(cookieName)
		if err == nil {
			cookieVal.Field(i).SetString(val)
		}
	}
	return nil
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
