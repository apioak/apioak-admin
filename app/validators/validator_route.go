package validators

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	routerPathPrefixMessages = map[string]string{
		utils.LocalEn: "%s must start with [/]",
		utils.LocalZh: "%s必须以[/]开始",
	}
	routerPathDefaultPathPrefixMessages = map[string]string{
		utils.LocalEn: "%s It is temporarily not allowed to start with the default routing path[/*]",
		utils.LocalZh: "%s暂不允许以默认路由路径[/*]开头",
	}
	routerRequestMethodOneOfMessages = map[string]string{
		utils.LocalEn: "%s must be one of [%s]",
		utils.LocalZh: "%s必须是[%s]中的一个",
	}
)

type ValidatorRouterAddUpdate struct {
	ServiceResID   string `json:"service_res_id" zh:"所属服务" en:"Belonging service" binding:"omitempty"`
	UpstreamResID  string `json:"upstream_res_id" zh:"上游服务" en:"Upstream service" binding:"omitempty"`
	RouterName     string `json:"router_name" zh:"路由名称" en:"Router name" binding:"omitempty"`
	RequestMethods string `json:"request_methods" zh:"请求方法" en:"Request method" binding:"required,min=3,CheckRouterRequestMethodOneOf"`
	RouterPath     string `json:"router_path" zh:"路由路径" en:"Routing path" binding:"required,min=2,CheckRouterPathPrefix"`
	Enable         int    `json:"enable" zh:"路由开关" en:"Routing enable" binding:"required,oneof=1 2"`
	UpstreamAddUpdate
}

type ValidatorRouterList struct {
	ServiceResID string `form:"service_res_id" json:"service_res_id" zh:"所属服务" en:"Belonging service" binding:"omitempty"`
	Search       string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	Enable       int    `form:"enable" json:"enable" zh:"路由开关" en:"Routing enable" binding:"omitempty,oneof=1 2"`
	Release      int    `form:"release" json:"release" zh:"发布状态" en:"Release status" binding:"omitempty,oneof=1 2 3"`
	BaseListPage
}

type RouterUpdateName struct {
	Name string `json:"name" zh:"路由名称" en:"Router name" binding:"required,min=1,max=30"`
}

type RouterSwitchEnable struct {
	Enable int `json:"enable" zh:"路由开关" en:"Router enable" binding:"required,oneof=1 2"`
}

func CheckRouterPathPrefix(fl validator.FieldLevel) bool {
	routePath := strings.TrimSpace(fl.Field().String())

	match := strings.Index(routePath, "/")
	if match != 0 {
		var errMsg string
		errMsg = fmt.Sprintf(routerPathPrefixMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName())
		packages.SetAllCustomizeValidatorErrMsgs("CheckRouterPathPrefix", errMsg)
		return false
	}

	matchDefaultPath := strings.Index(routePath, "/*")
	if matchDefaultPath == 0 {
		var errMsg string
		errMsg = fmt.Sprintf(routerPathDefaultPathPrefixMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName())
		packages.SetAllCustomizeValidatorErrMsgs("CheckRouterPathPrefix", errMsg)
		return false
	}

	return true
}

func CheckRouterRequestMethodOneOf(fl validator.FieldLevel) bool {
	requestMethods := strings.TrimSpace(fl.Field().String())

	requestMethodsSlice := strings.Split(requestMethods, ",")
	allRequestMethods := utils.AllRequestMethod()

	tmpRequestMethodsMap := make(map[string]byte)
	for _, allRequestMethod := range allRequestMethods {
		tmpRequestMethodsMap[allRequestMethod] = 0
	}

	filterAfterRequestMethods := make([]string, 0)
	for _, requestMethod := range requestMethodsSlice {
		requestMethodUpper := strings.ToUpper(requestMethod)
		if len(requestMethodUpper) == 0 {
			continue
		}

		if requestMethodUpper == utils.RequestMethodALL {
			filterAfterRequestMethods = []string{utils.RequestMethodALL}
			break
		}

		_, exist := tmpRequestMethodsMap[requestMethodUpper]
		if exist {
			filterAfterRequestMethods = append(filterAfterRequestMethods, requestMethodUpper)
		}
	}

	if len(filterAfterRequestMethods) == 0 {

		var errMsg string
		errMsg = fmt.Sprintf(routerRequestMethodOneOfMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName(), strings.Join(allRequestMethods, " "))
		packages.SetAllCustomizeValidatorErrMsgs("CheckRouterRequestMethodOneOf", errMsg)
		return false
	}

	return true
}

func GetRouterAttributesDefault(routerAddUpdate *ValidatorRouterAddUpdate) {
	routerAddUpdate.ServiceResID = strings.TrimSpace(routerAddUpdate.ServiceResID)
	routerAddUpdate.RouterPath = strings.TrimSpace(routerAddUpdate.RouterPath)
	routerAddUpdate.RequestMethods = strings.TrimSpace(routerAddUpdate.RequestMethods)
	routerAddUpdate.RouterName = strings.TrimSpace(routerAddUpdate.RouterName)

	requestMethodsSlice := strings.Split(routerAddUpdate.RequestMethods, ",")
	allRequestMethods := utils.AllRequestMethod()

	tmpRequestMethodsMap := make(map[string]byte)
	for _, allRequestMethod := range allRequestMethods {
		tmpRequestMethodsMap[allRequestMethod] = 0
	}

	filterAfterRequestMethods := make([]string, 0)
	requestMethodsMap := make(map[string]byte, 0)
	for _, requestMethod := range requestMethodsSlice {
		requestMethodUpper := strings.ToUpper(requestMethod)
		if len(requestMethodUpper) == 0 {
			continue
		}

		if requestMethodUpper == utils.RequestMethodALL {
			filterAfterRequestMethods = []string{utils.RequestMethodALL}
			break
		}

		_, exist := tmpRequestMethodsMap[requestMethodUpper]
		if !exist {
			continue
		}

		_, uinuqeExist := requestMethodsMap[requestMethodUpper]
		if uinuqeExist {
			continue
		}

		filterAfterRequestMethods = append(filterAfterRequestMethods, requestMethodUpper)
		requestMethodsMap[requestMethodUpper] = 0
	}

	if len(filterAfterRequestMethods) == (len(allRequestMethods) - 1) {
		filterAfterRequestMethods = []string{utils.RequestMethodALL}
	}

	routerAddUpdate.RequestMethods = strings.Join(filterAfterRequestMethods, ",")
}
