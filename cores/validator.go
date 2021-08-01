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

		// 注册获取的自定义 json tag方法（验证错误信息以传递参数名称为准）
		RegisterTag(validatorEngine, "json", "zh")

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(conf.Validator.Locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", conf.Validator.Locale)
		}

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
