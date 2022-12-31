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

type UpstreamTimeout struct {
	ReadTimeout    int `json:"read_timeout" zh:"读超市" en:"Read timeout" binding:"omitempty,min=1,max=600000"`
	WriteTimeout   int `json:"write_timeout" zh:"写超时" en:"Write timeout" binding:"omitempty,min=1,max=600000"`
	ConnectTimeout int `json:"connect_timeout" zh:"连接超时" en:"Connect timeout" binding:"omitempty,min=1,max=600000"`
}

type UpstreamAddUpdate struct {
	LoadBalance   int                     `json:"load_balance" zh:"负载均衡算法" en:"Load balancing algorithm" binding:"omitempty,CheckLoadBalanceOneOf"`
	UpstreamNodes []UpstreamNodeAddUpdate `json:"upstream_nodes" zh:"上游节点" en:"Upstream nodes" binding:"omitempty,CheckUpstreamNode"`
	UpstreamTimeout
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
