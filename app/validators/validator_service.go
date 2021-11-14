package validators

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

var (
	connectionTimeout             = 3000
	sendTimeout                   = 4000
	readTimeout                   = 5000
	connectionTimeoutKey          = "connection_timeout"
	sendTimeoutKey                = "send_timeout"
	readTimeoutKey                = "read_timeout"
	loadBalanceOneOfErrorMessages = map[string]string{
		utils.LocalEn: "%s must be one of [%s]",
		utils.LocalZh: "%s必须是[%s]中的一个",
	}
)

type ServiceAddUpdate struct {
	Timeouts       map[string]uint32      `json:"timeouts" zh:"超时时间" en:"Time out" binding:"omitempty"`
	LoadBalance    int                    `json:"load_balance" zh:"负载均衡算法" en:"Load balancing algorithm" binding:"omitempty,CheckLoadBalanceOneOf"`
	IsEnable       int                    `json:"is_enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	WebSocket      int                    `json:"web_socket" zh:"WebSocket" en:"WebSocket" binding:"omitempty,oneof=1 2"`
	HealthCheck    int                    `json:"health_check" zh:"健康检查" en:"Health" binding:"omitempty,oneof=1 2"`
	Protocol       int                    `json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	ServiceNodes   []ServiceNodeAddUpdate `json:"service_nodes" zh:"上游节点" en:"Service nodes" binding:"required,min=1,CheckServiceNode"`
	ServiceDomains []string               `json:"service_domains" zh:"域名" en:"Service domains" binding:"required,min=1,CheckServiceDomain"`
}

type ServiceList struct {
	Protocol int    `form:"protocol" json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	IsEnable int    `form:"is_enable" json:"is_enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	Search   string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

type ServiceUpdateName struct {
	Name string `json:"name" zh:"服务名称" en:"Service name" binding:"required,min=1,max=30"`
}

type ServiceSwitchEnable struct {
	IsEnable int `json:"is_enable" zh:"服务开关" en:"Service enable" binding:"required,oneof=1 2"`
}

type ServiceSwitchRelease struct {
	IsRelease int `json:"is_release" zh:"服务发布" en:"Service release" binding:"required,oneof=1"`
}

type ServiceSwitchWebsocket struct {
	WebSocket int `json:"web_socket" zh:"WebSocket" en:"WebSocket" binding:"required,oneof=1 2"`
}

type ServiceSwitchHealthCheck struct {
	HealthCheck int `json:"health_check" zh:"健康检查" en:"Health" binding:"required,oneof=1 2"`
}

func CheckLoadBalanceOneOf(fl validator.FieldLevel) bool {
	serviceLoadBalanceId := fl.Field().Int()

	loadBalanceInfos := utils.LoadBalanceList()

	loadBalanceIdsMap := make(map[int]byte, 0)
	loadBalanceIds := make([]string, 0)
	if len(loadBalanceInfos) != 0 {
		for _, loadBalanceInfo := range loadBalanceInfos {
			if loadBalanceInfo.Id == 0 {
				continue
			}

			loadBalanceIds = append(loadBalanceIds, strconv.Itoa(loadBalanceInfo.Id))
			loadBalanceIdsMap[loadBalanceInfo.Id] = 0
		}
	}

	_, exist := loadBalanceIdsMap[int(serviceLoadBalanceId)]
	if !exist {
		var errMsg string
		errMsg = fmt.Sprintf(loadBalanceOneOfErrorMessages[strings.ToLower(packages.GetValidatorLocale())], fl.FieldName(), strings.Join(loadBalanceIds, " "))
		packages.SetAllCustomizeValidatorErrMsgs("LoadBalanceOneOf", errMsg)
		return false
	}

	return true
}

func defaultServiceTimeOut() map[string]uint32 {
	timeInterface := make(map[string]uint32)

	timeInterface[connectionTimeoutKey] = uint32(connectionTimeout)
	timeInterface[sendTimeoutKey] = uint32(sendTimeout)
	timeInterface[readTimeoutKey] = uint32(readTimeout)

	return timeInterface
}

func CorrectServiceTimeOut(serviceTimeOuts *map[string]uint32) {
	defaultTimeOut := defaultServiceTimeOut()
	tmpServiceTimeOut := *serviceTimeOuts

	for timeOutKey, _ := range defaultTimeOut {
		timeOut, timeOutExist := tmpServiceTimeOut[timeOutKey]
		if timeOutExist {
			defaultTimeOut[timeOutKey] = timeOut
		}
	}

	serviceTimeOuts = &defaultTimeOut
}

func CorrectServiceAttributesDefault(serviceAddUpdate *ServiceAddUpdate) {
	if serviceAddUpdate.Protocol == 0 {
		serviceAddUpdate.Protocol = utils.ProtocolHTTP
	}
	if serviceAddUpdate.HealthCheck == 0 {
		serviceAddUpdate.HealthCheck = utils.EnableOff
	}
	if serviceAddUpdate.WebSocket == 0 {
		serviceAddUpdate.WebSocket = utils.EnableOff
	}
	if serviceAddUpdate.IsEnable == 0 {
		serviceAddUpdate.IsEnable = utils.EnableOff
	}
	if serviceAddUpdate.LoadBalance == 0 {
		serviceAddUpdate.LoadBalance = utils.LoadBalanceRoundRobin
	}
}
