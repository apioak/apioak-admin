package validators

type ClusterNodeList struct {
	NodeStatus int    `json:"node_status" zh:"节点状态" en:"Node status" binding:"omitempty,oneof=1 2"`
	IsEnable   int    `json:"is_enable" zh:"节点开关" en:"Node enable" binding:"omitempty,oneof=1 2"`
	Search     string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}
