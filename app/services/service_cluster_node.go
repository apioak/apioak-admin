package services

import (
	"apioak-admin/app/models"
	"apioak-admin/app/validators"
)

type ClusterNodeListInfo struct {
	ID         string `json:"id"`
	NodeIP     string `json:"node_ip"`
	NodeStatus int    `json:"node_status"`
	IsEnable   int    `json:"is_enable"`
}

func (c *ClusterNodeListInfo) ClusterNodeListPage(param *validators.ClusterNodeList) ([]ClusterNodeListInfo, int, error) {
	clusterNodesModel := models.ClusterNodes{}
	clusterNodeListInfos, total, clusterNodeListInfosErr := clusterNodesModel.ClusterNodeListPage(param)

	clusterNodeList := make([]ClusterNodeListInfo, 0)
	if len(clusterNodeListInfos) != 0 {
		for _, clusterNodeListInfo := range clusterNodeListInfos {
			clusterNodeInfo := ClusterNodeListInfo{}
			clusterNodeInfo.ID = clusterNodeListInfo.ID
			clusterNodeInfo.NodeIP = clusterNodeListInfo.NodeIP
			clusterNodeInfo.NodeStatus = clusterNodeListInfo.NodeStatus
			clusterNodeInfo.IsEnable = clusterNodeListInfo.IsEnable

			clusterNodeList = append(clusterNodeList, clusterNodeInfo)
		}
	}

	return clusterNodeList, total, clusterNodeListInfosErr
}
