package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/rpc"
	"apioak-admin/app/utils"
	"errors"
	"strings"
)

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
