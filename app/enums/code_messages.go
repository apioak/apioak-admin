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

	SwitchNoChange      = 101 // 开关无变化
	SwitchONProhibitsOp = 102 // 开关打开状态禁止该操作

	ServiceNull       = 10001 // 服务不存在
	ServiceParamsNull = 10002 // 服务参数缺失

	ServiceDomainExist = 10101 // 服务域名已存在

	RouteDefaultPathNoPermission = 10201 // [/*]默认路径暂无权限操作
	RoutePathExist               = 10202 // 路由路径已存在
	RouteNull                    = 10203 // 路由不存在
	RouteServiceNoMatch          = 10204 // 路由不在指定服务下
)

var ZhMapMessages = map[int]string{
	Success: "成功",
	Error:   "失败",

	SwitchNoChange:      "开关无变化",
	SwitchONProhibitsOp: "开关打开状态禁止该操作",

	ServiceNull:       "服务不存在",
	ServiceParamsNull: "服务参数缺失",

	ServiceDomainExist: "[%s]域名已存在",

	RouteDefaultPathNoPermission: "[/*]默认路径暂无权限操作",
	RoutePathExist:               "[%s]路由路径已存在",
	RouteNull:                    "路由不存在",
	RouteServiceNoMatch:          "路由不在指定服务下",
}

var EnMapMessages = map[int]string{
	Success: "success",
	Error:   "error",

	SwitchNoChange:      "No change in switch",
	SwitchONProhibitsOp: "This operation is prohibited when the switch is open",

	ServiceNull:       "Service does not exist",
	ServiceParamsNull: "Missing service parameters",

	ServiceDomainExist: "[%s]Domain name already exists",

	RouteDefaultPathNoPermission: "[/*]The default path does not have permission to operate temporarily",
	RoutePathExist:               "[%s]Routing path already exists",
	RouteNull:                    "Route does not exist",
	RouteServiceNoMatch:          "The route is not under the specified service",
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
