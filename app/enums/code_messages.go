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

	SwitchNoChange       = 101 // 开关无变化
	SwitchONProhibitsOp  = 102 // 开关打开状态禁止该操作
	ParamsError          = 103 // 参数异常
	EnablePublishedONOp  = 104 // 开关关闭状态并且已发布才允许该操作
	SwitchPublished      = 105 // 已发布
	IdConflict           = 106 // ID冲突
	ToReleaseProhibitsOp = 107 // 待发布状态禁止该操作

	ServiceNull       = 10001 // 服务不存在
	ServiceParamsNull = 10002 // 服务参数缺失

	ServiceDomainExist       = 10101 // [%s]服务域名已存在
	ServiceDomainFormatError = 10102 // 服务域名格式错误
	ServiceDomainSslNull     = 10104 // [%s]服务域名证书缺失

	RouteDefaultPathNoPermission    = 10201 // [/*]默认路径暂无权限操作
	RoutePathExist                  = 10202 // 路由路径已存在
	RouteNull                       = 10203 // 路由不存在
	RouteServiceNoMatch             = 10204 // 路由不在指定服务下
	RoutePluginExist                = 10205 // 路由插件已存在
	RoutePluginNull                 = 10206 // 路由插件不存在
	RoutePluginFormatError          = 10207 // 路由插件配置参数格式有误或参数错误
	RouteDefaultPathForbiddenPrefix = 10208 // [/*]默认路径暂不能作为路由开头
	RouteDefaultPathNull            = 10209 // [/*]默认路径路由不存在

	PluginTagExist   = 10301 // 插件标识已存在
	PluginNull       = 10302 // 插件不存在
	PluginRouteExist = 10303 // 插件已被路由绑定，暂不允许该操作
	PluginTagNull    = 10304 // 插件标识不存在，请检查[插件信息 - 插件标识]

	CertificateFormatError = 10401 // 证书格式错误
	CertificateParseError  = 10402 // 证书解析失败
	CertificateExist       = 10403 // 证书已存在
	CertificateNull        = 10404 // 证书不存在
	CertificateDomainExist = 10405 // 证书已被域名绑定，暂不支持该操作
	CertificateNoRelease   = 10406 // [%s]证书未发布
	CertificateEnableOff   = 10407 // [%s]证书未开启

	ClusterNodeNull = 10501 // 节点不存在

	UserEmailExist      = 10601 // 邮箱已注册
	UserNull            = 10602 // 用户不存在
	UserPasswordError   = 10603 // 用户与密码不匹配
	UserTokenError      = 10605 // 用户token验证信息失败
	UserNoLoggingIn     = 10606 // 用户未登录
	UserLoggingInError  = 10607 // 用户登录失败
	UserLoggingInExpire = 10608 // 用户登录已过期

	EtcdKeyNull     = 10701 // ETCD的key不能为空
	EtcdUnavailable = 10702 // ETCD不可用
)

var ZhMapMessages = map[int]string{
	Success: "成功",
	Error:   "失败",

	SwitchNoChange:       "开关无变化",
	SwitchONProhibitsOp:  "开关打开状态禁止该操作",
	ParamsError:          "参数异常",
	EnablePublishedONOp:  "开关关闭状态并且已发布才允许该操作",
	SwitchPublished:      "已发布",
	IdConflict:           "ID冲突",
	ToReleaseProhibitsOp: "待发布状态禁止该操作",

	ServiceNull:       "服务不存在",
	ServiceParamsNull: "服务参数缺失",

	ServiceDomainExist:       "[%s]域名已存在",
	ServiceDomainFormatError: "服务域名格式错误",
	ServiceDomainSslNull:     "[%s]服务域名证书缺失",

	RouteDefaultPathNoPermission:    "[/*]默认路径暂无权限操作",
	RoutePathExist:                  "[%s]路由路径已存在",
	RouteNull:                       "路由不存在",
	RouteServiceNoMatch:             "路由不在指定服务下",
	RoutePluginExist:                "路由插件已存在",
	RoutePluginNull:                 "路由插件不存在",
	RoutePluginFormatError:          "路由插件配置参数格式有误或参数错误",
	RouteDefaultPathForbiddenPrefix: "[/*]默认路径暂不能作为路由开头",
	RouteDefaultPathNull:            "[/*]默认路径路由不存在",

	PluginTagExist:   "插件标识已存在",
	PluginNull:       "插件不存在",
	PluginRouteExist: "插件已被路由绑定，暂不允许该操作",
	PluginTagNull:    "插件标识不存在，请检查[插件信息 - 插件标识]",

	CertificateFormatError: "证书格式错误",
	CertificateParseError:  "证书解析失败",
	CertificateExist:       "证书已存在",
	CertificateNull:        "证书不存在",
	CertificateDomainExist: "证书已被域名绑定，暂不允许该操作",
	CertificateNoRelease:   "[%s]证书未发布",
	CertificateEnableOff:   "[%s]证书未开启",

	ClusterNodeNull: "节点不存在",

	UserEmailExist:      "邮箱已注册",
	UserNull:            "用户不存在",
	UserPasswordError:   "用户与密码不匹配",
	UserTokenError:      "用户token验证信息失败",
	UserNoLoggingIn:     "用户未登录",
	UserLoggingInError:  "用户登录失败",
	UserLoggingInExpire: "用户登录已过期",

	EtcdKeyNull:     "ETCD的key不能为空",
	EtcdUnavailable: "ETCD不可用",
}

var EnMapMessages = map[int]string{
	Success: "success",
	Error:   "error",

	SwitchNoChange:       "No change in switch",
	SwitchONProhibitsOp:  "This operation is prohibited when the switch is open",
	ParamsError:          "Parameter abnormal",
	EnablePublishedONOp:  "The operation is only allowed when the switch is off and has been released",
	SwitchPublished:      "Published",
	IdConflict:           "ID conflict",
	ToReleaseProhibitsOp: "Pending status prohibits this operation",

	ServiceNull:       "Service does not exist",
	ServiceParamsNull: "Missing service parameters",

	ServiceDomainExist:       "[%s]Domain name already exists",
	ServiceDomainFormatError: "Service domain name format error",
	ServiceDomainSslNull:     "[%s]Service domain name certificate is missing",

	RouteDefaultPathNoPermission:    "[/*]The default path does not have permission to operate temporarily",
	RoutePathExist:                  "[%s]Routing path already exists",
	RouteNull:                       "Route does not exist",
	RouteServiceNoMatch:             "The route is not under the specified service",
	RoutePluginExist:                "Routing plugin already exists",
	RoutePluginNull:                 "The routing plugin does not exist",
	RoutePluginFormatError:          "The routing plugin configuration parameter format is incorrect or the parameter is wrong",
	RouteDefaultPathForbiddenPrefix: "[/*]The default path cannot be used as the beginning of the route",
	RouteDefaultPathNull:            "[/*]The default route route does not exist",

	PluginTagExist:   "Plugin tag already exists",
	PluginNull:       "Plugin does not exist",
	PluginRouteExist: "Plugin routing binding, operation is not allowed",
	PluginTagNull:    "The plugin tag does not exist, please check [plugin information - plugin tag]",

	CertificateFormatError: "Incorrect certificate format",
	CertificateParseError:  "Certificate parsing failed",
	CertificateExist:       "Certificate already exists",
	CertificateNull:        "Certificate does not exist",
	CertificateDomainExist: "The certificate has been bound by the domain name, operation is not allowed",
	CertificateNoRelease:   "[%s]Certificate not release",
	CertificateEnableOff:   "[%s]Certificate not enabled",

	ClusterNodeNull: "Node does not exist",

	UserEmailExist:      "Email has been registered",
	UserNull:            "User does not exist",
	UserPasswordError:   "User and password do not match",
	UserTokenError:      "User token verification information failed",
	UserNoLoggingIn:     "User is not logged in",
	UserLoggingInError:  "User login failed",
	UserLoggingInExpire: "User login has expired",

	EtcdKeyNull:     "The key of ETCD cannot be empty",
	EtcdUnavailable: "ETCD is not available",
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
