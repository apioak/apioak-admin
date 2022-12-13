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

func RouterUpstreamRelease(upstreamResIds []string, releaseType string) error {
	releaseType = strings.ToLower(releaseType)

	if (releaseType != utils.ReleaseTypePush) && (releaseType != utils.ReleaseTypeDelete) {
		return errors.New(enums.CodeMessages(enums.ReleaseTypeError))
	}

	upstreamModel := models.Upstreams{}

	upstreamList, err := upstreamModel.UpstreamListByResIds(upstreamResIds)
	if err != nil {
		return err
	}

	if len(upstreamList) == 0 {
		return nil
	}

	upstreamConfigList := make([]rpc.UpstreamConfig, 0)
	upstreamNodeResIds := make([]string, 0)
	for _, upstreamInfo := range upstreamList {

		upstreamConfig, upstreamConfigErr := generateUpstreamConfig(upstreamInfo)
		if upstreamConfigErr != nil {
			return upstreamConfigErr
		}

		if len(upstreamConfig.Name) == 0 {
			continue
		}

		upstreamConfigList = append(upstreamConfigList, upstreamConfig)

		if len(upstreamConfig.Nodes) > 0 {
			for _, upstreamCOnfigNodeName := range upstreamConfig.Nodes {
				upstreamNodeResIds = append(upstreamNodeResIds, upstreamCOnfigNodeName.Name)
			}
		}

	}

	newApiOak := rpc.NewApiOak()

	if releaseType == utils.ReleaseTypePush {
		releaseUpstreamNodeErr := UpstreamNodeRelease(upstreamNodeResIds, releaseType)
		if releaseUpstreamNodeErr != nil {
			return releaseUpstreamNodeErr
		}

		upstreamPutErr := newApiOak.UpstreamPut(upstreamConfigList)
		if upstreamPutErr != nil {
			return upstreamPutErr
		}
	} else {
		upstreamDeleteErr := newApiOak.UpstreamDelete(upstreamConfigList)
		if upstreamDeleteErr != nil {
			return upstreamDeleteErr
		}

		releaseUpstreamNodeErr := UpstreamNodeRelease(upstreamNodeResIds, releaseType)
		if releaseUpstreamNodeErr != nil {
			return releaseUpstreamNodeErr
		}
	}

	return nil
}

func generateUpstreamConfig(upstreamInfo models.Upstreams) (rpc.UpstreamConfig, error) {
	upstreamConfig := rpc.UpstreamConfig{}

	configBalanceList := utils.ConfigBalanceList()
	configBalanceMap := make(map[int]string)
	for _, configBalanceInfo := range configBalanceList {
		configBalanceMap[configBalanceInfo.Id] = configBalanceInfo.Name
	}

	upstreamConfig.Algorithm = utils.ConfigBalanceNameRoundRobin
	configBalance, ok := configBalanceMap[upstreamInfo.Algorithm]
	if ok {
		upstreamConfig.Algorithm = configBalance
	}

	upstreamConfig.Name = upstreamInfo.ResID
	upstreamConfig.ConnectTimeout = upstreamInfo.ConnectTimeout
	upstreamConfig.WriteTimeout = upstreamInfo.WriteTimeout
	upstreamConfig.ReadTimeout = upstreamInfo.ReadTimeout
	upstreamConfig.Nodes = make([]rpc.ConfigObjectName, 0)

	upstreamNodeModel := models.UpstreamNodes{}
	upstreamNodeList, err := upstreamNodeModel.UpstreamNodeListByUpstreamResIds([]string{upstreamInfo.ResID})
	if err != nil {
		return upstreamConfig, err
	}

	if len(upstreamNodeList) != 0 {
		for _, upstreamNodeInfo := range upstreamNodeList {
			upstreamConfig.Nodes = append(upstreamConfig.Nodes, rpc.ConfigObjectName{
				Name: upstreamNodeInfo.ResID,
			})
		}
	}

	return upstreamConfig, nil
}
