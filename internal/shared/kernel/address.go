package kernel

import (
	"gochat/internal/shared/errors"
	"strings"
	"unicode/utf8"
)

type Address struct {
	Country  string
	Province string
	City     string
	District string
	Street   string
}

func (a Address) String() string {
	var parts []string
	if a.Country != "" {
		parts = append(parts, a.Country)
	}
	if a.Province != "" {
		parts = append(parts, a.Province)
	}
	if a.City != "" {
		parts = append(parts, a.City)
	}
	if a.District != "" {
		parts = append(parts, a.District)
	}
	if a.Street != "" {
		parts = append(parts, a.Street)
	}

	return strings.Join(parts, " ")
}

// Validate 验证地址格式
func (a Address) Validate() error {
	// 验证完整地址长度
	fullAddress := a.String()

	// 检查地址是否为空
	if utf8.RuneCountInString(fullAddress) == 0 {
		return errors.ErrInvalidFormat
	}

	// 检查地址长度限制 (假设最大500字符)
	if utf8.RuneCountInString(fullAddress) > 500 {
		return errors.ErrInvalidLength
	}

	// 验证省份/城市信息
	if a.Province != "" && utf8.RuneCountInString(a.Province) > 50 {
		return errors.ErrInvalidLength
	}

	if a.City != "" && utf8.RuneCountInString(a.City) > 50 {
		return errors.ErrInvalidLength
	}

	return nil
}
