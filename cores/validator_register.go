package cores

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func RegisterTag(validatorEngine *validator.Validate, tag string) {
	validatorEngine.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get(tag), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func RegisterCustomizeTrans(validatorEngine *validator.Validate, trans ut.Translator, validatorRule string, customizeTrans string) {
	_ = validatorEngine.RegisterTranslation(validatorRule, trans, func(ut ut.Translator) error {
		return ut.Add(validatorRule, customizeTrans, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(validatorRule, fe.Field())
		return t
	})
}

