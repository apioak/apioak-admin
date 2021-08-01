package models

import "time"

type ClusterNodes struct {
	ID         string    `gorm:"column:id;primary_key"` //Cluster node id
	NodeIP     string    `gorm:"column:node_ip"`        //Node IP
	IPType     int       `gorm:"column:ip_type"`        //IP Type  0:IPV4  1:IPV6
	NodeStatus int       `gorm:"column:node_status"`    //Node status  1:health  2:Unhealthy
	IsEnable   int       `gorm:"column:is_enable"`      //Node enable  0:off  1:on
	CreatedAt  time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt  time.Time `gorm:"column:updated_at"`     //Update time
}

// TableName sets the insert table name for this struct type
func (c *ClusterNodes) TableName() string {
	return "oak_cluster_nodes"
}
