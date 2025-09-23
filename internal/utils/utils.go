package utils

import (
	"context"
	"errors"
	"gochat/internal/types"
	"gochat/internal/utils/timeout"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strings"
)

// NormalizeVarArgs 标准化可变参数，根据参数数量返回不同形式的结果
func NormalizeVarArgs[T any](data ...T) any {
	if l := len(data); l == 1 {
		return data[0]
	} else if l > 1 {
		return data
	} else {
		return nil
	}
}

// IsEmptyData 判断 data 是否为 nil 或空结构体（所有字段不可导出或被 json:"-" 标记）
func IsEmptyData(data any) bool {
	if data == nil {
		return true
	}
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Struct:
		// 空结构体
		return v.NumField() == 0
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	default:
		log.Fatalf("unhandled default case")
	}
	return false
}

// GetOption 获取可选的参数
func GetOption[T any](opts ...T) (bool, T) {
	if len(opts) == 0 {
		var zero T
		return false, zero
	} else {
		return true, opts[0]
	}
}

// HandleDatabaseError 统一将底层 DB/Redis/上下文错误映射为业务错误类型
func HandleDatabaseError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err == nil {
		return nil
	}

	// 1. 超时或取消优先
	if timeout.IsCanceledOrTimeout(err) {
		return types.ErrTimeout
	}

	// 2. gorm 未找到 -> 转为 types.ErrNotFound
	if e := CheckNotFoundError(err); errors.Is(e, types.ErrNotFound) {
		return types.ErrNotFound
	}

	// 3. 重复键 -> 转为 types.ErrDuplicateKey
	if e := CheckDuplicateKeyError(err); errors.Is(e, types.ErrDuplicateKey) {
		return types.ErrDuplicateKey
	}

	// 4. 其他错误保持原样返回，便于上层记录或透传
	return err
}

func CheckDuplicateKeyError(err error) error {
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "1062") {
		return types.ErrDuplicateKey
	}
	return err
}

func CheckNotFoundError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return types.ErrNotFound
	}
	return err
}

// HasStructField 检查结构体是否包含特定字段
func HasStructField(obj interface{}, fieldName string) bool {
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

// GetFieldAddr 获取结构体中指定字段的地址
func GetFieldAddr(obj interface{}, fieldName string) interface{} {
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

func IsNotFound(err error) bool {
	return errors.Is(err, types.ErrNotFound)
}

func IsDuplicate(err error) bool {
	return errors.Is(err, types.ErrDuplicateKey)
}
