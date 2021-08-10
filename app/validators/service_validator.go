package validators

type ServiceAdd struct {
	Timeouts     string `json:"timeouts" zh:"超时时间" en:"Time out" binding:"omitempty,json"`
	LoadBalance  int    `json:"load_balance" zh:"负载均衡算法" en:"Load balancing algorithm" binding:"omitempty,oneof=1 2"`
	IsEnable     int    `json:"is_enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	WebSocket    int    `json:"web_socket" zh:"WebSocket" en:"WebSocket" binding:"omitempty,oneof=1 2"`
	HealthCheck  int    `json:"health_check" zh:"健康检查" en:"Health" binding:"omitempty,oneof=1 2"`
	Protocol     int    `json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	ServiceNodes string `json:"service_nodes" zh:"上游节点" en:"Service nodes" binding:"required,json,CheckServiceNode"`
	ServiceNames string `json:"service_domains" zh:"域名" en:"Service domains" binding:"required,CheckServiceDomain"`
}
