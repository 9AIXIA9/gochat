package utils

import (
	"context"
	"errors"
	"fmt"
	"gochat/internal/types"
	"gorm.io/gorm"
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

// IsCanceledOrTimeout 错误是上下文被取消或者是超时
func IsCanceledOrTimeout(err error) bool {
	return IsCanceled(err) || IsTimeout(err)
}

// IsCanceled 判断错误是否由上下文取消引起
func IsCanceled(err error) bool {
	return errors.Is(err, types.ErrCanceled) || errors.Is(err, context.Canceled)
}

// IsTimeout 判断错误是否由上下文超时引起
func IsTimeout(err error) bool {
	return errors.Is(err, types.ErrTimeout) || errors.Is(err, context.DeadlineExceeded)
}

// IsEmptyData 判断 data 是否为 nil 或空结构体（所有字段不可导出或被 json:"-" 标记）
func IsEmptyData(data any) bool {
	if data == nil {
		return true
	}

	v := reflect.ValueOf(data)

	// 解包指针类型
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		if t.NumField() == 0 {
			return true
		}
		// 遍历字段，若存在可导出且未被 json:"-" 标记的字段，则认为不是空
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			// 非导出字段跳过
			if f.PkgPath != "" {
				continue
			}
			// 被 json:"-" 标记的字段跳过
			if f.Tag.Get("json") == "-" {
				continue
			}
			// 找到可序列化的字段，认为不是空
			return false
		}
		// 所有字段都不可序列化
		return true
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		// 其他类型不视为“空结构体”，让其正常序列化
		return false
	}
}

// HandleDatabaseError 统一将底层 db/Redis/上下文错误映射为业务错误类型
func HandleDatabaseError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err == nil {
		return nil
	}

	// 1. 超时或取消优先
	if IsCanceledOrTimeout(err) {
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

// ValidateAllSubStructsNotEmpty 递归验证传入结构体的所有子结构体/切片/映射等是否为空。
// - obj 必须为结构体或指向结构体的指针。
// - 可通过 struct tag `validate:"skip"` 跳过某个字段的校验。
// 返回非 nil 错误时，错误信息列出所有被判空的字段路径。
func ValidateAllSubStructsNotEmpty(obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("object is nil")
	}

	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return fmt.Errorf("object is nil pointer")
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct or pointer to struct")
	}

	var empties []string

	var walk func(reflect.Value, string)
	walk = func(val reflect.Value, path string) {
		if !val.IsValid() {
			return
		}

		// 解包指针/接口
		for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
			if val.IsNil() {
				empties = append(empties, path)
				return
			}
			val = val.Elem()
		}

		switch val.Kind() {
		case reflect.Struct:
			t := val.Type()
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				// 跳过未导出字段
				if f.PkgPath != "" {
					continue
				}
				// 支持 tag 跳过
				if f.Tag.Get("validate") == "skip" {
					continue
				}

				fv := val.Field(i)
				fieldPath := f.Name
				if path != "" {
					fieldPath = path + "." + f.Name
				}

				// 使用已有的 IsEmptyData 判空规则
				// 需要注意：fv.Interface() 在导出字段上是安全的
				if IsEmptyData(fv.Interface()) {
					empties = append(empties, fieldPath)
					// 为空时无需继续深入该字段
					continue
				}

				// 若是容器或嵌套结构，继续递归检查其内部元素
				switch fv.Kind() {
				case reflect.Struct, reflect.Ptr, reflect.Interface:
					walk(fv, fieldPath)
				case reflect.Slice, reflect.Array:
					if fv.Len() > 0 {
						for j := 0; j < fv.Len(); j++ {
							walk(fv.Index(j), fmt.Sprintf("%s[%d]", fieldPath, j))
						}
					}
				case reflect.Map:
					if fv.Len() > 0 {
						for _, k := range fv.MapKeys() {
							walk(fv.MapIndex(k), fmt.Sprintf("%s[%v]", fieldPath, k.Interface()))
						}
					}
				default:
				}
			}
		case reflect.Slice, reflect.Array:
			if val.Len() == 0 {
				empties = append(empties, path)
			} else {
				for i := 0; i < val.Len(); i++ {
					walk(val.Index(i), fmt.Sprintf("%s[%d]", path, i))
				}
			}
		case reflect.Map:
			if val.Len() == 0 {
				empties = append(empties, path)
			} else {
				for _, k := range val.MapKeys() {
					walk(val.MapIndex(k), fmt.Sprintf("%s[%v]", path, k.Interface()))
				}
			}
		default:
			// 基础类型不用递归
		}
	}

	walk(v, "")

	if len(empties) > 0 {
		return fmt.Errorf("empty fields: %s", strings.Join(empties, ", "))
	}
	return nil
}
