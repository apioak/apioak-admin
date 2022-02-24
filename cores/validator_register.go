package cores

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/validators"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func RegisterTag(validatorEngine *validator.Validate, tag string, translation string) {

	tag = strings.TrimSpace(tag)
	translation = strings.TrimSpace(translation)
	if translation == "" {
		translation = "en"
	}

	validatorEngine.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get(tag), ",", -1)[0]
		translationStr := strings.SplitN(field.Tag.Get(translation), ",", -1)[0]
		if name == "-" {
			return ""
		}
		return translationStr + "[" + name + "]"
	})
}

func RegisterCustomizeValidator(validatorEngine *validator.Validate) error {
	if err := validatorEngine.RegisterValidation("CheckServiceDomain", validators.CheckServiceDomain); err != nil {
		return err
	}
	if err := validatorEngine.RegisterValidation("CheckServiceNode", validators.CheckServiceNode); err != nil {
		return err
	}
	if err := validatorEngine.RegisterValidation("CheckLoadBalanceOneOf", validators.CheckLoadBalanceOneOf); err != nil {
		return err
	}

	if err := validatorEngine.RegisterValidation("CheckRoutePathPrefix", validators.CheckRoutePathPrefix); err != nil {
		return err
	}
	if err := validatorEngine.RegisterValidation("CheckRouteRequestMethodOneOf", validators.CheckRouteRequestMethodOneOf); err != nil {
		return err
	}

	if err := validatorEngine.RegisterValidation("CheckPluginTypeOneOf", validators.CheckPluginTypeOneOf); err != nil {
		return err
	}
	if err := validatorEngine.RegisterValidation("CheckPluginTagOneOf", validators.CheckPluginTagOneOf); err != nil {
		return err
	}

	packages.SetCustomizeValidator(validatorEngine)
	return nil
}