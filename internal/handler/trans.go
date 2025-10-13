package handler

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"reflect"
	"strings"
)

var (
	globalValidator *validator.Validate
	globalTrans     ut.Translator
)

// GetTrans 获取全局翻译器
func GetTrans() ut.Translator {
	return globalTrans
}

// InitTrans 初始化验证器和翻译器
func InitTrans(locale string) error {
	// 创建验证器实例
	globalValidator = validator.New()

	// 配置验证器
	// 注册JSON标签名函数
	globalValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// 根据配置选择语言
	var uni *ut.UniversalTranslator
	var err error

	switch locale {
	case "en":
		// 创建英语翻译器
		enTranslator := en.New()
		uni = ut.New(enTranslator)
		globalTrans, _ = uni.GetTranslator(enTranslator.Locale())

		// 注册英语翻译
		err = enTranslations.RegisterDefaultTranslations(globalValidator, globalTrans)
		if err != nil {
			return fmt.Errorf("register English translation failed, err:%w", err)
		}

		// 注册英语自定义翻译
		if err := registerEnCustomTranslations(globalValidator, globalTrans); err != nil {
			return err
		}

	default: // 默认使用中文
		// 创建中文翻译器
		zhTranslator := zh.New()
		uni = ut.New(zhTranslator)
		globalTrans, _ = uni.GetTranslator(zhTranslator.Locale())

		// 注册中文翻译
		err = zhTranslations.RegisterDefaultTranslations(globalValidator, globalTrans)
		if err != nil {
			return fmt.Errorf("register Chinese translation failed, err:%w", err)
		}

		// 注册中文自定义翻译
		if err := registerZhCustomTranslations(globalValidator, globalTrans); err != nil {
			return err
		}
	}

	log.Printf("translator initialized with locale: %s", locale)
	return nil
}

// registerZhCustomTranslations 注册中文自定义翻译
func registerZhCustomTranslations(v *validator.Validate, trans ut.Translator) error {
	translations := []struct {
		tag         string
		translation string
		override    bool
	}{
		{
			tag:         "required",
			translation: "{0}不能为空",
			override:    true,
		},
		{
			tag:         "min",
			translation: "{0}长度必须至少为{1}个字符",
			override:    true,
		},
		{
			tag:         "max",
			translation: "{0}长度不能超过{1}个字符",
			override:    true,
		},
	}

	return registerTranslations(v, trans, translations)
}

// registerEnCustomTranslations 注册英文自定义翻译
func registerEnCustomTranslations(v *validator.Validate, trans ut.Translator) error {
	translations := []struct {
		tag         string
		translation string
		override    bool
	}{
		{
			tag:         "required",
			translation: "{0} cannot be empty",
			override:    true,
		},
		{
			tag:         "min",
			translation: "{0} must be at least {1} characters long",
			override:    true,
		},
		{
			tag:         "max",
			translation: "{0} cannot exceed {1} characters",
			override:    true,
		},
	}

	return registerTranslations(v, trans, translations)
}

// registerTranslations 注册翻译
func registerTranslations(v *validator.Validate, trans ut.Translator, translations []struct {
	tag         string
	translation string
	override    bool
}) error {
	for _, t := range translations {
		err := v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		if err != nil {
			return fmt.Errorf("register translation failed, err:%v", err)
		}
	}
	return nil
}

// registrationFunc 返回一个注册函数
func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

// translateFunc 翻译函数
func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		return fe.Error() // 回退到默认错误
	}
	return t
}
