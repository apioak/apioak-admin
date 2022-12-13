package models

type PluginConfigs struct {
	ID          int    `gorm:"column:id;primary_key"` // primary key
	ResID       string `gorm:"column:res_id"`         // Plugin config id
	Name        string `gorm:"column:name"`           // Plugin config name
	Type        int    `gorm:"column:type"`           // Plugin relation type 1:service  2:router
	TargetID    string `gorm:"column:target_id"`      // Target id
	PluginResID string `gorm:"column:plugin_res_id"`  // Plugin res id
	PluginKey   string `gorm:"column:plugin_key"`     // Plugin key
	Config      string `gorm:"column:config"`         // Plugin configuration
	Enable      int    `gorm:"column:enable"`         // Plugin config enable  1:on  2:off
	ModelTime
}

// TableName sets the insert table name for this struct type
func (p *PluginConfigs) TableName() string {
	return "oak_plugin_configs"
}
