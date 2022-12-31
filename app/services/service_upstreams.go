package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/rpc"
	"apioak-admin/app/utils"
	"errors"
	"strings"
)

type UpstreamItem struct {
	ResID          string `json:"res_id"`
	Algorithm      int    `json:"algorithm"`
	ConnectTimeout int    `json:"connect_timeout"`
	WriteTimeout   int    `json:"write_timeout"`
	ReadTimeout    int    `json:"read_timeout"`
}

func (u UpstreamItem) UpstreamDetailByResId(resId string) (upstreamItem UpstreamItem, err error) {

	upstreamModel := models.Upstreams{}
	upstreamDetail, err := upstreamModel.UpstreamDetailByResId(resId)
	if err != nil {
		return
	}

	if len(upstreamDetail.ResID) == 0 {
		return
	}

	upstreamItem.ResID = upstreamDetail.ResID
	upstreamItem.Algorithm = upstreamDetail.Algorithm
	upstreamItem.ConnectTimeout = upstreamDetail.ConnectTimeout
	upstreamItem.WriteTimeout = upstreamDetail.WriteTimeout
	upstreamItem.ReadTimeout = upstreamDetail.ReadTimeout

	return
}

func UpstreamRelease(upstreamResIds []string, releaseType string) (err error) {
	if len(upstreamResIds) == 0 {
		return
	}

	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		err = errors.New(enums.CodeMessages(enums.ReleaseTypeError))
		return
	}

	newApiOak := rpc.NewApiOak()
	if releaseType == utils.ReleaseTypePush {
		getUpstreamConfigList := make([]rpc.UpstreamConfig, 0)
		getUpstreamConfigList, err = newApiOak.UpstreamGet(upstreamResIds)

		if err != nil {
			return
		}

		upstreamNodeIds := make([]string, 0)
		for _, getUpstreamConfigInfo := range getUpstreamConfigList {

			if len(getUpstreamConfigInfo.Nodes) == 0 {
				continue
			}

			for _, nodeInfo := range getUpstreamConfigInfo.Nodes {
				upstreamNodeIds = append(upstreamNodeIds, nodeInfo.Id)
			}
		}

		err = newApiOak.UpstreamDelete(upstreamResIds)
		if err != nil {
			return
		}

		err = newApiOak.UpstreamNodeDeleteByIds(upstreamNodeIds)
		if err != nil {
			return
		}

		err = UpstreamNodeRelease(upstreamResIds, releaseType)
		if err != nil {
			return
		}

		upstreamModel := models.Upstreams{}
		var upstreamList []models.Upstreams
		upstreamList, err = upstreamModel.UpstreamListByResIds(upstreamResIds)

		if err != nil {
			return
		}

		if len(upstreamList) == 0 {
			return
		}

		upstreamNodeModel := models.UpstreamNodes{}
		upstreamNodeList := make([]models.UpstreamNodes, 0)
		upstreamNodeList, err = upstreamNodeModel.UpstreamNodeListByUpstreamResIds(upstreamResIds)
		if err != nil {
			return
		}

		upstreamNodeListMap := make(map[string]models.UpstreamNodes)
		for _, upstreamNodeInfo := range upstreamNodeList {
			upstreamNodeListMap[upstreamNodeInfo.UpstreamResID] = upstreamNodeInfo
		}

		upstreamConfigList := make([]rpc.UpstreamConfig, 0)
		for _, upstreamInfo := range upstreamList {
			_, ok := upstreamNodeListMap[upstreamInfo.ResID]
			if !ok {
				continue
			}

			var upstreamConfig rpc.UpstreamConfig
			upstreamConfig, err = generateUpstreamConfig(upstreamInfo)
			if err != nil {
				return
			}

			upstreamConfigList = append(upstreamConfigList, upstreamConfig)
		}

		err = newApiOak.UpstreamPut(upstreamConfigList)
		if err != nil {
			return
		}

	} else {
		err = newApiOak.UpstreamDelete(upstreamResIds)
		if err != nil {
			return
		}

		err = UpstreamNodeRelease(upstreamResIds, releaseType)
		if err != nil {
			return
		}
	}

	return
}

func generateUpstreamConfig(upstreamInfo models.Upstreams) (config rpc.UpstreamConfig, err error) {

	configBalanceList := utils.ConfigBalanceList()
	configBalanceMap := make(map[int]string)
	for _, configBalanceInfo := range configBalanceList {
		configBalanceMap[configBalanceInfo.Id] = configBalanceInfo.Name
	}

	config.Algorithm = utils.ConfigBalanceNameRoundRobin
	configBalance, ok := configBalanceMap[upstreamInfo.Algorithm]
	if ok {
		config.Algorithm = configBalance
	}

	config.Name = upstreamInfo.ResID
	config.ConnectTimeout = upstreamInfo.ConnectTimeout
	config.WriteTimeout = upstreamInfo.WriteTimeout
	config.ReadTimeout = upstreamInfo.ReadTimeout
	config.Nodes = make([]rpc.ConfigObjectName, 0)

	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList := make([]models.UpstreamNodes, 0)
	upstreamNodeList, err = upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{upstreamInfo.ResID})
	if err != nil {
		return
	}

	if len(upstreamNodeList) != 0 {
		for _, upstreamNodeInfo := range upstreamNodeList {
			config.Nodes = append(config.Nodes, rpc.ConfigObjectName{
				Name: upstreamNodeInfo.ResID,
			})
		}
	}

	return
}
