package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"strings"
)

type ClusterNodes struct {
	ID         string `gorm:"column:id;primary_key"` //Cluster node id
	NodeIP     string `gorm:"column:node_ip"`        //Node IP
	IPType     int    `gorm:"column:ip_type"`        //IP Type  1:IPV4  2:IPV6
	NodeStatus int    `gorm:"column:node_status"`    //Node status  1:health  2:Unhealthy
	IsEnable   int    `gorm:"column:is_enable"`      //Node enable  1:on  2:off
	ModelTime
}

// TableName sets the insert table name for this struct type
func (c *ClusterNodes) TableName() string {
	return "oak_cluster_nodes"
}

var clusterNodeId = ""

func (c *ClusterNodes) ClusterNodeIdUnique(cIds map[string]string) (string, error) {
	if c.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeClusterNode)
		if err != nil {
			return "", err
		}
		c.ID = tmpID
	}

	result := packages.GetDb().
		Table(c.TableName()).
		Select("id").
		First(&c)

	mapId := cIds[c.ID]
	if (result.RowsAffected == 0) && (c.ID != mapId) {
		clusterNodeId = c.ID
		cIds[c.ID] = c.ID
		return clusterNodeId, nil
	} else {
		nodeId, nodeIdErr := utils.IdGenerate(utils.IdTypeClusterNode)
		if nodeIdErr != nil {
			return "", nodeIdErr
		}
		c.ID = nodeId
		_, err := c.ClusterNodeIdUnique(cIds)
		if err != nil {
			return "", err
		}
	}

	return clusterNodeId, nil
}

func (c *ClusterNodes) ClusterNodeInfoById(id string) ClusterNodes {
	clusterNodeInfo := ClusterNodes{}
	packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		First(&clusterNodeInfo)

	return clusterNodeInfo
}

func (c *ClusterNodes) ClusterNodeSwitchEnable(id string, enable int) error {
	updateErr := packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		Update("is_enable", enable).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (c *ClusterNodes) ClusterNodeListPage(param *validators.ClusterNodeList) (list []ClusterNodes, total int, listError error) {
	tx := packages.GetDb().
		Table(c.TableName())

	if param.IsEnable != 0 {
		tx = tx.Where("is_enable = ?", param.IsEnable)
	}
	if param.NodeStatus != 0 {
		tx = tx.Where("node_status = ?", param.NodeStatus)
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where("node_ip LIKE ?", search)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.Order("updated_at DESC")
	listError = ListPaginate(tx, &list, &param.BaseListPage)
	return

}

func (c *ClusterNodes) ClusterNodeDelete(id string) error {
	deleteErr := packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		Delete(c).Error

	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (c *ClusterNodes) ClusterNodeInfoByIp(ip string) ClusterNodes {
	clusterNodeInfo := ClusterNodes{}
	db := packages.GetDb().
		Table(c.TableName()).
		Where("node_ip = ?", ip)

	db.First(&clusterNodeInfo)

	return clusterNodeInfo
}

func (c *ClusterNodes) ClusterNodeAdd(clusterNodesData *ClusterNodes) error {
	tpmIds := map[string]string{}
	clusterNodeIdUnique, clusterNodeIdUniqueErr := c.ClusterNodeIdUnique(tpmIds)
	if clusterNodeIdUniqueErr != nil {
		return clusterNodeIdUniqueErr
	}
	clusterNodesData.ID = clusterNodeIdUnique

	err := packages.GetDb().
		Table(c.TableName()).
		Create(clusterNodesData).Error

	return err
}

func (c *ClusterNodes) ClusterNodeUpdate(id string, clusterNodesData *ClusterNodes) error {
	updateError := packages.GetDb().
		Table(c.TableName()).
		Where("id = ?", id).
		Updates(clusterNodesData).Error

	return updateError
}
