package validators

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
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
	Timeouts       string `json:"timeouts" zh:"超时时间" en:"Time out" binding:"omitempty,json"`
	LoadBalance    int    `json:"load_balance" zh:"负载均衡算法" en:"Load balancing algorithm" binding:"omitempty,LoadBalanceOneOf"`
	IsEnable       int    `json:"is_enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	WebSocket      int    `json:"web_socket" zh:"WebSocket" en:"WebSocket" binding:"omitempty,oneof=1 2"`
	HealthCheck    int    `json:"health_check" zh:"健康检查" en:"Health" binding:"omitempty,oneof=1 2"`
	Protocol       int    `json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	ServiceNodes   string `json:"service_nodes" zh:"上游节点" en:"Service nodes" binding:"required,json,CheckServiceNode"`
	ServiceDomains string `json:"service_domains" zh:"域名" en:"Service domains" binding:"required,CheckServiceDomain"`
}

type ServiceList struct {
	Protocol int    `form:"protocol" zh:"请求协议" json:"protocol" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	IsEnable int    `form:"is_enable" zh:"服务开关" json:"is_enable" en:"Service enable" binding:"omitempty,oneof=1 2"`
	Search   string `form:"search" zh:"搜索内容" json:"search" en:"Search content" binding:"omitempty"`
	BaseListPage
}

type ServiceUpdateName struct {
	Name string `json:"name" zh:"服务名称" en:"Service name" binding:"required,min=1,max=30"`
}

type ServiceSwitchEnable struct {
	IsEnable int `json:"is_enable" zh:"服务开关" en:"Service enable" binding:"required,oneof=1 2"`
}

type ServiceSwitchWebsocket struct {
	WebSocket int `json:"web_socket" zh:"WebSocket" en:"WebSocket" binding:"required,oneof=1 2"`
}

type ServiceSwitchHealthCheck struct {
	HealthCheck int `json:"health_check" zh:"健康检查" en:"Health" binding:"required,oneof=1 2"`
}

func LoadBalanceOneOf(fl validator.FieldLevel) bool {
	serviceLoadBalanceId := fl.Field().Int()

	loadBalance := utils.LoadBalance{}
	loadBalanceInfos := loadBalance.LoadBalanceList()

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

func GetServiceAddTimeOut(times string) string {

	defaultTimeOut := defaultServiceTimeOut()
	if len(times) <= 0 {

		timeStr, err := json.Marshal(defaultTimeOut)
		if err != nil {
			return ""
		}
		return string(timeStr)
	}

	timeInterface := make(map[string]interface{})
	jsonErr := json.Unmarshal([]byte(times), &timeInterface)
	if jsonErr != nil {
		timeStr, err := json.Marshal(defaultTimeOut)
		if err != nil {
			return ""
		}
		return string(timeStr)
	}

	for timeKey, millisecond := range timeInterface {
		switch timeKey {
		case connectionTimeoutKey:

			millisecondInt, err := strconv.Atoi(millisecond.(string))
			if err != nil {
				break
			}
			defaultTimeOut[connectionTimeoutKey] = uint32(millisecondInt)
		case sendTimeoutKey:
			millisecondInt, err := strconv.Atoi(millisecond.(string))
			if err != nil {
				break
			}
			defaultTimeOut[sendTimeoutKey] = uint32(millisecondInt)
		case readTimeoutKey:
			millisecondInt, err := strconv.Atoi(millisecond.(string))
			if err != nil {
				break
			}
			defaultTimeOut[readTimeoutKey] = uint32(millisecondInt)
		}
	}

	timeStr, err := json.Marshal(defaultTimeOut)
	if err != nil {
		return ""
	}
	return string(timeStr)
}

func GetServiceAttributesDefault(serviceInfo ServiceAddUpdate) ServiceAddUpdate {
	if serviceInfo.Protocol == 0 {
		serviceInfo.Protocol = utils.ProtocolHTTP
	}
	if serviceInfo.HealthCheck == 0 {
		serviceInfo.HealthCheck = utils.EnableOff
	}
	if serviceInfo.WebSocket == 0 {
		serviceInfo.WebSocket = utils.EnableOff
	}
	if serviceInfo.IsEnable == 0 {
		serviceInfo.IsEnable = utils.EnableOff
	}
	if serviceInfo.LoadBalance == 0 {
		serviceInfo.LoadBalance = utils.LoadBalanceRoundRobin
	}

	return serviceInfo
}
