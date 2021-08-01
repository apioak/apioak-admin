package cores

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func RegisterTag(validatorEngine *validator.Validate, tag string, zh string) {
	validatorEngine.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get(tag), ",", -1)[0]
		zhStr := strings.SplitN(field.Tag.Get(zh), ",", -1)[0]
		if name == "-" {
			return ""
		}
		return zhStr + "[" + name + "]"
	})
}


