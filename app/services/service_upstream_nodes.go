package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/rpc"
	"apioak-admin/app/utils"
	"errors"
	"strings"
)

type UpstreamNodeItem struct {
	ResID         string `json:"res_id"`
	UpstreamResID string `json:"upstream_res_id"`
	NodeIP        string `json:"node_ip"`
	IPType        int    `json:"ip_type"`
	IPTypeName    string `json:"ip_type_name"`
	NodePort      int    `json:"node_port"`
	NodeWeight    int    `json:"node_weight"`
	Health        int    `json:"health"`
	HealthName    string `json:"health_name"`
	HealthCheck   int    `json:"health_check"`
}

func (n UpstreamNodeItem) UpstreamNodeListByUpstreamResIds(upstreamResIds []string) (nodeList []UpstreamNodeItem, err error) {
	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByUpstreamResIds(upstreamResIds)
	if err != nil || len(upstreamNodeList) == 0 {
		return
	}

	iPTypeNameMap := utils.IpIdNameMap()
	healthTypeNameMap := utils.HealthTypeNameMap()

	for _, upstreamNodeDetail := range upstreamNodeList {

		nodeList = append(nodeList, UpstreamNodeItem{
			ResID: upstreamNodeDetail.ResID,
			UpstreamResID: upstreamNodeDetail.UpstreamResID,
			NodeIP: upstreamNodeDetail.NodeIP,
			IPType: upstreamNodeDetail.IPType,
			IPTypeName: iPTypeNameMap[upstreamNodeDetail.IPType],
			NodePort: upstreamNodeDetail.NodePort,
			NodeWeight: upstreamNodeDetail.NodeWeight,
			Health: upstreamNodeDetail.Health,
			HealthName: healthTypeNameMap[upstreamNodeDetail.Health],
		})
	}

	return
}

func UpstreamNodeRelease(upstreamNodeResIds []string, releaseType string) error {
	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		return errors.New(enums.CodeMessages(enums.ReleaseTypeError))
	}

	upstreamNodeModel := models.UpstreamNodes{}

	upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByResIds(upstreamNodeResIds)
	if err != nil {
		return err
	}

	if len(upstreamNodeList) == 0 {
		return nil
	}

	upstreamNodeConfigList := make([]rpc.UpstreamNodeConfig, 0)
	for _, upstreamNodeInfo := range upstreamNodeList {

		upstreamNodeConfig, upstreamNodeConfigErr := generateUpstreamNodeConfig(upstreamNodeInfo)
		if upstreamNodeConfigErr != nil {
			return upstreamNodeConfigErr
		}

		if len(upstreamNodeConfig.Name) == 0 {
			continue
		}

		upstreamNodeConfigList = append(upstreamNodeConfigList, upstreamNodeConfig)
	}

	newApiOak := rpc.NewApiOak()

	if releaseType == utils.ReleaseTypePush {
		upstreamNodePutErr := newApiOak.UpstreamNodePut(upstreamNodeConfigList)
		if upstreamNodePutErr != nil {
			return upstreamNodePutErr
		}
	} else {
		upstreamNodeDeleteErr := newApiOak.UpstreamNodeDelete(upstreamNodeConfigList)
		if upstreamNodeDeleteErr != nil {
			return upstreamNodeDeleteErr
		}
	}

	return err
}

func generateUpstreamNodeConfig(upstreamNodeInfo models.UpstreamNodes) (rpc.UpstreamNodeConfig, error) {
	upstreamNodeConfig := rpc.UpstreamNodeConfig{}

	configHealthList := utils.ConfigUpstreamNodeHealthList()
	configHealthMap := make(map[int]string)
	for _, configHealthInfo := range configHealthList {
		configHealthMap[configHealthInfo.Id] = configHealthInfo.Name
	}

	upstreamNodeConfig.Health = utils.ConfigHealthY
	configHealth, ok := configHealthMap[upstreamNodeInfo.Health]
	if ok {
		upstreamNodeConfig.Health = configHealth
	}

	upstreamNodeConfig.Name = upstreamNodeInfo.ResID
	upstreamNodeConfig.Address = upstreamNodeInfo.NodeIP
	upstreamNodeConfig.Port = upstreamNodeInfo.NodePort
	upstreamNodeConfig.Weight = upstreamNodeInfo.NodeWeight

	// @todo 节点健康检查
	upstreamNodeConfig.Check.Enabled = false

	return upstreamNodeConfig, nil
}
