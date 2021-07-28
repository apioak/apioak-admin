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

	trans ut.Translator
)

func InitValidator(conf*ConfigGlobal) (err error) {
	if validatorEngine, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册获取 json 的自定义方法
		RegisterTag(validatorEngine, "json")
		RegisterTag(validatorEngine, "form")

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(conf.Validator.Locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", conf.Validator.Locale)
		}

		// 增加额外的翻译
		RegisterCustomizeTrans(validatorEngine, trans, "required_with", "{0} 为必填字段!")
		RegisterCustomizeTrans(validatorEngine, trans, "required_without", "{0} 为必填字段!")
		RegisterCustomizeTrans(validatorEngine, trans, "required_without_all", "{0} 为必填字段!")

		// 注册翻译器
		switch strings.ToLower(conf.Validator.Locale) {
		case LocaleEn:
			err = enTranslations.RegisterDefaultTranslations(validatorEngine, trans)
		case LocaleZh:
			err = zhTranslations.RegisterDefaultTranslations(validatorEngine, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(validatorEngine, trans)
		}
		packages.SetTranslator(&trans)
		return
	}
	return
}
