package validators

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	certificateContentRequiredMessages = map[string]string{
		utils.LocalEn: "[content] in %sJSON content is a required field",
		utils.LocalZh: "%sJSON内容中[content]为必填字段",
	}
)

type CertificateAddUpdate struct {
	Certificate string `json:"certificate" zh:"证书内容" en:"Certificate content" binding:"required,json,CheckCertificateContentRequired"`
	PrivateKey  string `json:"private_key" zh:"私钥内容" en:"Private key content" binding:"required,json,CheckCertificateContentRequired"`
	IsEnable    int    `json:"is_enable" zh:"证书开关" en:"Certificate enable" binding:"required,oneof=1 2"`
}

type CertificateList struct {
	IsEnable int    `form:"is_enable" json:"is_enable" zh:"证书开关" en:"Certificate enable" binding:"omitempty,oneof=1 2"`
	Search   string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

type CertificateSwitchEnable struct {
	IsEnable int `json:"is_enable" zh:"证书开关" en:"Certificate enable" binding:"required,oneof=1 2"`
}

type CertificateSwitchRelease struct {
	IsRelease int `json:"is_release" zh:"证书发布" en:"Certificate release" binding:"required,oneof=1"`
}

func CheckCertificateContentRequired(fl validator.FieldLevel) bool {
	certificateContent := strings.TrimSpace(fl.Field().String())

	type contentStruct struct {
		Content string `json:"content"`
	}

	contentInfo := contentStruct{}
	contentInfoErr := json.Unmarshal([]byte(certificateContent), &contentInfo)
	if contentInfoErr != nil {
		packages.SetAllCustomizeValidatorErrMsgs("CheckCertificateContentRequired", fmt.Sprintf("json.Unmarshal error: ["+contentInfoErr.Error()+"]"))
		return false
	}

	if len(contentInfo.Content) == 0 {
		errMsg := fmt.Sprintf(certificateContentRequiredMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName())
		packages.SetAllCustomizeValidatorErrMsgs("CheckCertificateContentRequired", errMsg)
		return false
	}

	return true
}
