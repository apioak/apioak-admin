package models

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
