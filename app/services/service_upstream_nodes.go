package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/rpc"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"strconv"
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
	nodeList = make([]UpstreamNodeItem, 0)
	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByUpstreamResIds(upstreamResIds)
	if err != nil || len(upstreamNodeList) == 0 {
		return
	}

	iPTypeNameMap := utils.IpIdNameMap()
	healthTypeNameMap := utils.HealthTypeNameMap()

	for _, upstreamNodeDetail := range upstreamNodeList {

		nodeList = append(nodeList, UpstreamNodeItem{
			ResID:         upstreamNodeDetail.ResID,
			UpstreamResID: upstreamNodeDetail.UpstreamResID,
			NodeIP:        upstreamNodeDetail.NodeIP,
			IPType:        upstreamNodeDetail.IPType,
			IPTypeName:    iPTypeNameMap[upstreamNodeDetail.IPType],
			NodePort:      upstreamNodeDetail.NodePort,
			NodeWeight:    upstreamNodeDetail.NodeWeight,
			Health:        upstreamNodeDetail.Health,
			HealthName:    healthTypeNameMap[upstreamNodeDetail.Health],
			HealthCheck:   upstreamNodeDetail.HealthCheck,
		})
	}

	return
}

func DiffUpstreamNode(upstreamResID string, paramNodeList []validators.UpstreamNodeAddUpdate) (
	addNodeList []models.UpstreamNodes, updateNodeList []models.UpstreamNodes, delNodeResIds []string) {

	if len(upstreamResID) == 0 {
		return
	}

	paramNodeListMap := make(map[string]validators.UpstreamNodeAddUpdate)
	for _, paramNodeInfo := range paramNodeList {
		paramNodeListMapKey := paramNodeInfo.NodeIp + "-" + strconv.Itoa(paramNodeInfo.NodePort)
		paramNodeListMap[paramNodeListMapKey] = paramNodeInfo
	}

	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{upstreamResID})
	if err != nil {
		return
	}

	upstreamNodeListMap := make(map[string]models.UpstreamNodes)
	for _, upstreamNodeInfo := range upstreamNodeList {
		upstreamNodeListMapKey := upstreamNodeInfo.NodeIP + "-" + strconv.Itoa(upstreamNodeInfo.NodePort)
		upstreamNodeListMap[upstreamNodeListMapKey] = upstreamNodeInfo

		paramNodeInfo, ok := paramNodeListMap[upstreamNodeListMapKey]

		if ok {
			updateNodeList = append(updateNodeList, models.UpstreamNodes{
				ResID:      upstreamNodeInfo.ResID,
				NodePort:   paramNodeInfo.NodePort,
				NodeWeight: paramNodeInfo.NodeWeight,
				Health:     paramNodeInfo.Health,
			})
		} else {
			delNodeResIds = append(delNodeResIds, upstreamNodeInfo.ResID)
		}
	}

	ipNameIdMap := utils.IpNameIdMap()

	for _, paramNodeListInfo := range paramNodeList {
		upstreamNodeListMapKey := paramNodeListInfo.NodeIp + "-" + strconv.Itoa(paramNodeListInfo.NodePort)
		_, ok := upstreamNodeListMap[upstreamNodeListMapKey]
		if !ok {

			resId, resIdErr := upstreamNodeModel.ModelUniqueId()
			if resIdErr != nil {
				continue
			}

			ipType, ipTypeErr := utils.DiscernIP(paramNodeListInfo.NodeIp)
			if ipTypeErr != nil {
				continue
			}

			addNodeList = append(addNodeList, models.UpstreamNodes{
				ResID:         resId,
				UpstreamResID: upstreamResID,
				NodeIP:        paramNodeListInfo.NodeIp,
				IPType:        ipNameIdMap[ipType],
				NodePort:      paramNodeListInfo.NodePort,
				NodeWeight:    paramNodeListInfo.NodeWeight,
				Health:        paramNodeListInfo.Health,
				HealthCheck:   utils.HealthCheckOff,
			})
		}
	}

	return
}

func UpstreamNodeLocalCloudDiff(localNodeList []models.UpstreamNodes, cloudNodeList []rpc.UpstreamNodeConfig) (
	putNodeIds []string, deleteNodeIds []string) {

	localNodeListMap := make(map[string]models.UpstreamNodes)
	for _, localNodeInfo := range localNodeList {
		localNodeListMap[localNodeInfo.ResID] = localNodeInfo
		putNodeIds = append(putNodeIds, localNodeInfo.ResID)
	}

	for _, cloudNodeInfo := range cloudNodeList {
		if _, exits := localNodeListMap[cloudNodeInfo.Name]; !exits {
			deleteNodeIds = append(deleteNodeIds, cloudNodeInfo.Name)
		}
	}

	return
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
	upstreamNodeConfig.Check.Enabled = false

	return upstreamNodeConfig, nil
}

func NodeRelease(nodeResIds []string, releaseType string) (err error) {
	if len(nodeResIds) == 0 {
		return
	}

	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		err = errors.New(enums.CodeMessages(enums.ReleaseTypeError))
		return
	}

	newApiOak := rpc.NewApiOak()

	if releaseType == utils.ReleaseTypeDelete {
		err = newApiOak.UpstreamNodeDelete(nodeResIds)
		return
	} else {
		upstreamNodeModel := models.UpstreamNodes{}

		upstreamNodeList := make([]models.UpstreamNodes, 0)
		upstreamNodeList, err = upstreamNodeModel.UpstreamNodeListByResIds(nodeResIds)
		if err != nil {
			return
		}

		if len(upstreamNodeList) == 0 {
			return
		}

		upstreamNodeConfigList := make([]rpc.UpstreamNodeConfig, 0)
		for _, upstreamNodeInfo := range upstreamNodeList {
			var upstreamNodeConfig rpc.UpstreamNodeConfig
			upstreamNodeConfig, err = generateUpstreamNodeConfig(upstreamNodeInfo)
			if err != nil {
				return err
			}

			if len(upstreamNodeConfig.Name) == 0 {
				continue
			}

			upstreamNodeConfigList = append(upstreamNodeConfigList, upstreamNodeConfig)
		}

		err = newApiOak.UpstreamNodePut(upstreamNodeConfigList)
		if err != nil {
			return
		}
	}

	return
}
