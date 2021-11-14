package validators

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	routePathPrefixMessages = map[string]string{
		utils.LocalEn: "%s must start with [/]",
		utils.LocalZh: "%s必须以[/]开始",
	}
	routePathDefaultPathPrefixMessages = map[string]string{
		utils.LocalEn: "%s It is temporarily not allowed to start with the default routing path[/*]",
		utils.LocalZh: "%s暂不允许以默认路由路径[/*]开头",
	}
	routeRequestMethodOneOfMessages = map[string]string{
		utils.LocalEn: "%s must be one of [%s]",
		utils.LocalZh: "%s必须是[%s]中的一个",
	}
)

type ValidatorRouteAddUpdate struct {
	IsEnable       int    `json:"is_enable" zh:"路由开关" en:"Routing enable" binding:"required,oneof=1 2"`
	RouteName      string `json:"route_name" zh:"路由名称" en:"Route name" binding:"omitempty"`
	RequestMethods string `json:"request_methods" zh:"请求方法" en:"Request method" binding:"required,min=3,CheckRouteRequestMethodOneOf"`
	RoutePath      string `json:"route_path" zh:"路由路径" en:"Routing path" binding:"required,min=2,CheckRoutePathPrefix"`
	ServiceID      string `json:"service_id" zh:"服务ID" en:"Service id" binding:"omitempty"`
}

type ValidatorRouteList struct {
	Search   string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	IsEnable int    `form:"is_enable" json:"is_enable" zh:"路由开关" en:"Routing enable" binding:"omitempty,oneof=1 2"`
	BaseListPage
}

type RouteUpdateName struct {
	Name string `json:"name" zh:"路由名称" en:"Route name" binding:"required,min=1,max=30"`
}

type RouteSwitchEnable struct {
	IsEnable int `json:"is_enable" zh:"路由开关" en:"Route enable" binding:"required,oneof=1 2"`
}

func CheckRoutePathPrefix(fl validator.FieldLevel) bool {
	routePath := strings.TrimSpace(fl.Field().String())

	match := strings.Index(routePath, "/")
	if match != 0 {
		var errMsg string
		errMsg = fmt.Sprintf(routePathPrefixMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName())
		packages.SetAllCustomizeValidatorErrMsgs("CheckRoutePathPrefix", errMsg)
		return false
	}

	matchDefaultPath := strings.Index(routePath, "/*")
	if matchDefaultPath == 0 {
		var errMsg string
		errMsg = fmt.Sprintf(routePathDefaultPathPrefixMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName())
		packages.SetAllCustomizeValidatorErrMsgs("CheckRoutePathPrefix", errMsg)
		return false
	}

	return true
}

func CheckRouteRequestMethodOneOf(fl validator.FieldLevel) bool {
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
		errMsg = fmt.Sprintf(routeRequestMethodOneOfMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName(), strings.Join(allRequestMethods, " "))
		packages.SetAllCustomizeValidatorErrMsgs("CheckRouteRequestMethodOneOf", errMsg)
		return false
	}

	return true
}

func GetRouteAttributesDefault(routeAddUpdate *ValidatorRouteAddUpdate) {
	routeAddUpdate.ServiceID = strings.TrimSpace(routeAddUpdate.ServiceID)
	routeAddUpdate.RoutePath = strings.TrimSpace(routeAddUpdate.RoutePath)
	routeAddUpdate.RequestMethods = strings.TrimSpace(routeAddUpdate.RequestMethods)
	routeAddUpdate.RouteName = strings.TrimSpace(routeAddUpdate.RouteName)

	requestMethodsSlice := strings.Split(routeAddUpdate.RequestMethods, ",")
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

	routeAddUpdate.RequestMethods = strings.Join(filterAfterRequestMethods, ",")
}
