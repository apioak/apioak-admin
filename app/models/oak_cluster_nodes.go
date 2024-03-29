package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"strings"
)

type ClusterNodes struct {
	ID         string `gorm:"column:id;primary_key"` //Cluster node id
	NodeIP     string `gorm:"column:node_ip"`        //Node IP
	IPType     int    `gorm:"column:ip_type"`        //IP Type  1:IPV4  2:IPV6
	NodeStatus int    `gorm:"column:node_status"`    //Node status  1:health  2:Unhealthy
	ModelTime
}

// TableName sets the insert table name for this struct type
func (c *ClusterNodes) TableName() string {
	return "oak_cluster_nodes"
}

var recursionTimesClusterNodes = 1

func (m *ClusterNodes) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeClusterNode)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("id = ?", generateId).
		Select("id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesClusterNodes = 1
		return generateId, nil
	} else {
		if recursionTimesClusterNodes == utils.IdGenerateMaxTimes {
			recursionTimesClusterNodes = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesClusterNodes++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
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

	if param.IPType != 0 {
		tx = tx.Where("ip_type = ?", param.IPType)
	}
	if param.NodeStatus != 0 {
		tx = tx.Where("node_status = ?", param.NodeStatus)
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		tx = tx.Where(packages.GetDb().
			Table(c.TableName()).
			Where("node_ip LIKE ?", search).
			Or("id LIKE ?", search))
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.Order("created_at DESC")
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
	clusterNodeIdUnique, clusterNodeIdUniqueErr := c.ModelUniqueId()
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
