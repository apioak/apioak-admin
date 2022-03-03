package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
)

func CheckClusterNodeNull(id string) error {
	clusterNodesModel := models.ClusterNodes{}
	clusterNodeInfo := clusterNodesModel.ClusterNodeInfoById(id)
	if clusterNodeInfo.ID != id {
		return errors.New(enums.CodeMessages(enums.ClusterNodeNull))
	}

	return nil
}

func CheckClusterNodeExist(ip string) error {
	clusterNodesModel := models.ClusterNodes{}
	clusterNodeInfo := clusterNodesModel.ClusterNodeInfoByIp(ip)
	if len(clusterNodeInfo.ID) != 0 {
		return errors.New(enums.CodeMessages(enums.ClusterNodeExist))
	}

	return nil
}

func ClusterNodeAdd(clusterNodeAdd *validators.ClusterNodeAdd) error {
	ipTypeName, ipTypeNameErr := utils.DiscernIP(clusterNodeAdd.NodeIP)
	if ipTypeNameErr != nil {
		return ipTypeNameErr
	}

	ipType, ipTypeErr := utils.IPNameToType(ipTypeName)
	if ipTypeErr != nil {
		return ipTypeErr
	}

	clusterNodesModel := models.ClusterNodes{
		NodeIP:     clusterNodeAdd.NodeIP,
		IPType:     ipType,
		NodeStatus: clusterNodeAdd.NodeStatus,
	}

	addErr := clusterNodesModel.ClusterNodeAdd(&clusterNodesModel)
	if addErr != nil {
		return addErr
	}

	return nil
}

type ClusterNodeListInfo struct {
	ID         string `json:"id"`
	NodeIP     string `json:"node_ip"`
	IPType     int    `json:"ip_type"`
	NodeStatus int    `json:"node_status"`
}

func (c *ClusterNodeListInfo) ClusterNodeListPage(param *validators.ClusterNodeList) ([]ClusterNodeListInfo, int, error) {
	clusterNodesModel := models.ClusterNodes{}
	clusterNodeListInfos, total, clusterNodeListInfosErr := clusterNodesModel.ClusterNodeListPage(param)

	clusterNodeList := make([]ClusterNodeListInfo, 0)
	if len(clusterNodeListInfos) != 0 {
		for _, clusterNodeListInfo := range clusterNodeListInfos {
			clusterNodeInfo := ClusterNodeListInfo{}
			clusterNodeInfo.ID = clusterNodeListInfo.ID
			clusterNodeInfo.IPType = clusterNodeListInfo.IPType
			clusterNodeInfo.NodeIP = clusterNodeListInfo.NodeIP
			clusterNodeInfo.NodeStatus = clusterNodeListInfo.NodeStatus

			clusterNodeList = append(clusterNodeList, clusterNodeInfo)
		}
	}

	return clusterNodeList, total, clusterNodeListInfosErr
}

func ClusterNodeDelete(id string) error {
	clusterNodesModel := models.ClusterNodes{}
	deleteErr := clusterNodesModel.ClusterNodeDelete(id)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
