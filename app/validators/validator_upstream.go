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
	defaultTimeout                = 3000
	loadBalanceOneOfErrorMessages = map[string]string{
		utils.LocalEn: "%s must be one of [%s]",
		utils.LocalZh: "%s必须是[%s]中的一个",
	}
)

type UpstreamList struct {
	Enable    int    `form:"enable" json:"enable" zh:"上游开关" en:"Upstream enable" binding:"omitempty,oneof=1 2"`
	Release   int    `form:"release" json:"release" zh:"发布状态" en:"Release status" binding:"omitempty,oneof=1 2 3"`
	Algorithm int    `form:"algorithm" json:"algorithm" zh:"负载均衡" en:"Load balancing" binding:"omitempty"`
	Search    string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

type UpstreamTimeout struct {
	ReadTimeout    int `json:"read_timeout" zh:"读超时" en:"Read timeout" binding:"omitempty,min=1,max=600000"`
	WriteTimeout   int `json:"write_timeout" zh:"写超时" en:"Write timeout" binding:"omitempty,min=1,max=600000"`
	ConnectTimeout int `json:"connect_timeout" zh:"连接超时" en:"Connect timeout" binding:"omitempty,min=1,max=600000"`
}

type UpstreamAddUpdate struct {
	Name        string `json:"name" zh:"上游名称" en:"Upstream name" binding:"omitempty,min=1,max=30"`
	LoadBalance int    `json:"load_balance" zh:"负载均衡算法" en:"Load balancing algorithm" binding:"omitempty,CheckLoadBalanceOneOf"`
	Enable      int    `json:"enable" zh:"上游开关" en:"Upstream enable" binding:"omitempty,oneof=1 2"`
	UpstreamTimeout
	UpstreamNodes []UpstreamNodeAddUpdate `json:"upstream_nodes" zh:"上游节点" en:"Upstream nodes" binding:"required,min=1,CheckUpstreamNode"`
}

type UpstreamUpdateName struct {
	Name string `json:"name" zh:"上游名称" en:"Upstream name" binding:"required,min=1,max=30"`
}

type UpstreamSwitchEnable struct {
	Enable int `json:"enable" zh:"上游开关" en:"Upstream enable" binding:"required,oneof=1 2"`
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

func CorrectUpstreamDefault(upstreamData *UpstreamAddUpdate) {
	if upstreamData.LoadBalance == 0 {
		upstreamData.LoadBalance = utils.LoadBalanceRoundRobin
	}
	if upstreamData.ConnectTimeout == 0 {
		upstreamData.ConnectTimeout = defaultTimeout
	}
	if upstreamData.WriteTimeout == 0 {
		upstreamData.WriteTimeout = defaultTimeout
	}
	if upstreamData.ReadTimeout == 0 {
		upstreamData.ReadTimeout = defaultTimeout
	}
}
