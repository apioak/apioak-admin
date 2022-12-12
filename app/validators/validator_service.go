package validators

import (
	"apioak-admin/app/utils"
)

type ServiceAddUpdate struct {
	Name           string   `json:"name" zh:"服务名称" en:"Service name" binding:"required,min=1,max=30"`
	Enable         int      `json:"enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	Release        int      `json:"release" zh:"发布开关" en:"Release status enable" binding:"omitempty,oneof=1 2"`
	Protocol       int      `json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	ServiceDomains []string `json:"service_domains" zh:"域名" en:"Service domains" binding:"required,min=1,CheckServiceDomain"`
}

type ServiceList struct {
	Protocol int    `form:"protocol" json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	Enable   int    `form:"enable" json:"enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	Release  int    `form:"release" json:"release" zh:"发布状态" en:"Release status" binding:"omitempty,oneof=1 2 3"`
	Search   string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}

type ServiceUpdateName struct {
	Name string `json:"name" zh:"服务名称" en:"Service name" binding:"required,min=1,max=30"`
}

type ServiceSwitchEnable struct {
	IsEnable int `json:"is_enable" zh:"服务开关" en:"Service enable" binding:"required,oneof=1 2"`
}

type ServiceSwitchWebsocket struct {
	WebSocket int `json:"web_socket" zh:"WebSocket" en:"WebSocket" binding:"required,oneof=1 2"`
}

type ServiceSwitchHealthCheck struct {
	HealthCheck int `json:"health_check" zh:"健康检查" en:"Health" binding:"required,oneof=1 2"`
}

func CorrectServiceAttributesDefault(serviceAddUpdate *ServiceAddUpdate) {
	if serviceAddUpdate.Protocol == 0 {
		serviceAddUpdate.Protocol = utils.ProtocolHTTP
	}
	if serviceAddUpdate.Enable == 0 {
		serviceAddUpdate.Enable = utils.EnableOff
	}
}
