package enums

import (
	"apioak-admin/app/packages"
	"strings"
)

var (
	localEn = "en"
	localZh = "zh"
)

const (
	Success = 0  // 成功
	Error   = -1 // 失败

	SwitchNoChange = 101 // 开关无变化

	ServiceNull       = 10001 // 服务不存在
	ServiceParamsNull = 10002 // 服务不存在

	ServiceDomainExist = 10101 // 服务域名已存在
)

var ZhMapMessages = map[int]string{
	Success: "成功",
	Error:   "失败",

	SwitchNoChange: "开关无变化",

	ServiceNull:       "服务不存在",
	ServiceParamsNull: "服务参数缺失",

	ServiceDomainExist: "[%s]域名已存在",
}

var EnMapMessages = map[int]string{
	Success: "success",
	Error:   "error",

	SwitchNoChange: "No change in switch",

	ServiceNull:       "Service does not exist",
	ServiceParamsNull: "Missing service parameters",

	ServiceDomainExist: "[%s]Domain name already exists",
}

func CodeMessages(code int) string {
	mapMessages := EnMapMessages
	if getLocal() == localZh {
		mapMessages = ZhMapMessages
	}
	return mapMessages[code]
}

func getLocal() string {
	globalLocal := packages.GetValidatorLocale()
	var local = localEn
	if strings.ToLower(globalLocal) == localZh {
		local = localZh
	}
	return local
}
