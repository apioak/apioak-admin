package validators

type RouterPluginAddUpdate struct {
	Config   interface{} `json:"config" zh:"路由插件配置" en:"Routing plugin configuration" binding:"required"`
	Order    int         `json:"order" zh:"插件执行顺序" en:"Plugin execution order" binding:"required,min=1,max=30"`
	IsEnable int         `json:"is_enable" zh:"路由插件开关" en:"Routing plugin enable" binding:"required,oneof=1 2"`
	PluginID string      `json:"plugin_id" binding:"omitempty"`
	RouterID  string      `json:"router_id" binding:"omitempty"`
}

type RouterPluginSwitchEnable struct {
	IsEnable int `json:"is_enable" zh:"路由插件开关" en:"Routing plugin enable" binding:"required,oneof=1 2"`
}
