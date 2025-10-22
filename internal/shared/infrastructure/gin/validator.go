package gin

import (
	"context"
	"errors"
	"fmt"
	httputils "gochat/internal/shared/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewValidator() (*Validator, error) {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	enT := en.New()
	uni := ut.New(enT)
	translator, _ := uni.GetTranslator(enT.Locale())

	if err := enTranslations.RegisterDefaultTranslations(v, translator); err != nil {
		return nil, fmt.Errorf("register English translation failed, err: %w", err)
	}
	if err := registerEnCustomTranslations(v, translator); err != nil {
		return nil, err
	}

	return &Validator{
		validator:  v,
		translator: translator,
	}, nil
}

func (b *Validator) Validator() *validator.Validate { return b.validator }
func (b *Validator) Translator() ut.Translator      { return b.translator }

func (b *Validator) Validate(ctx context.Context, model any) (*httputils.ApiResponse, error) {
	if err := b.validator.StructCtx(ctx, model); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			message := b.buildValidationErrorMessage(validationErrors)
			return httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, message), nil
		}
		return nil, err
	}
	return nil, nil
}

func (b *Validator) buildValidationErrorMessage(typeErr validator.ValidationErrors) string {
	var sb strings.Builder
	for i, fe := range typeErr {
		if i > 0 {
			sb.WriteString("; ")
		}
		if t, err := b.translator.T(fe.Tag(), fe.Field(), fe.Param()); err == nil {
			sb.WriteString(fe.Field() + ": " + t)
		} else {
			sb.WriteString(fe.Field() + ": " + fe.Error())
		}
	}
	return sb.String()
}

func registerEnCustomTranslations(v *validator.Validate, trans ut.Translator) error {
	translations := []struct {
		tag         string
		translation string
		override    bool
	}{
		{tag: "required", translation: "{0} cannot be empty", override: true},
		{tag: "min", translation: "{0} must be at least {1} characters long", override: true},
		{tag: "max", translation: "{0} cannot exceed {1} characters", override: true},
	}
	return registerTranslations(v, trans, translations)
}

func registerTranslations(v *validator.Validate, trans ut.Translator, translations []struct {
	tag         string
	translation string
	override    bool
}) error {
	for _, t := range translations {
		if err := v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc); err != nil {
			return fmt.Errorf("register translation failed, err:%v", err)
		}
	}
	return nil
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(tr ut.Translator) error {
		return tr.Add(tag, translation, override)
	}
}

func translateFunc(tr ut.Translator, fe validator.FieldError) string {
	t, err := tr.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		return fe.Error()
	}
	return t
}
