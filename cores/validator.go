package cores

import (
	"apioak-admin/app/packages"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)

var (
	LocaleEn = "en"
	LocaleZh = "zh"
	trans    ut.Translator
)

func InitValidator(conf *ConfigGlobal) (err error) {

	if validatorEngine, ok := binding.Validator.Engine().(*validator.Validate); ok {

		uni := ut.New(en.New(), zh.New(), en.New())
		confLocal := strings.ToLower(conf.Validator.Locale)
		trans, ok = uni.GetTranslator(confLocal)
		if !ok {
			return fmt.Errorf("uni.GetTranslator (%s) failed", confLocal)
		}

		registerValidatorErr := RegisterCustomizeValidator(validatorEngine)
		if registerValidatorErr != nil {
			return fmt.Errorf("Custom registration verification error (%s)", registerValidatorErr)
		}

		RegisterTag(validatorEngine, "form", confLocal)
		RegisterTag(validatorEngine, "json", confLocal)

		switch confLocal {
		case LocaleEn:
			err = enTranslations.RegisterDefaultTranslations(validatorEngine, trans)
		case LocaleZh:
			err = zhTranslations.RegisterDefaultTranslations(validatorEngine, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(validatorEngine, trans)
		}
		packages.SetValidatorLocale(confLocal)
		packages.SetTranslator(&trans)
		return
	}
	return
}

