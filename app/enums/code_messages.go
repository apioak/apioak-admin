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

	ServiceNull = 10001 // 服务不存在

	ServiceDomainExist = 10101 // 服务域名已存在
)

var ZhMapMessages = map[int]string{
	Success: "成功",
	Error:   "失败",

	ServiceNull: "服务不存在",

	ServiceDomainExist: "[%s]域名已存在",
}

var EnMapMessages = map[int]string{
	Success: "success",
	Error:   "error",

	ServiceNull: "Service does not exist",

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
