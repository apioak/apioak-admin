package validators

type ClusterNodeAdd struct {
	NodeIP     string `form:"node_ip" json:"node_ip" zh:"节点IP" en:"Node IP" binding:"required,ip"`
	NodeStatus int    `form:"node_status" json:"node_status" zh:"节点健康状态" en:"Node health status" binding:"omitempty,oneof=1 2"`
}

type ClusterNodeList struct {
	IPType     int    `form:"ip_type" json:"ip_type" zh:"IP类型" en:"IP type" binding:"omitempty,oneof=1 2"`
	NodeStatus int    `form:"node_status" json:"node_status" zh:"节点状态" en:"Node status" binding:"omitempty,oneof=1 2"`
	Search     string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}
