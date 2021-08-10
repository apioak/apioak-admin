package models

type RoutePlugins struct {
	ID       string `gorm:"column:id;primary_key"` //Plugin id
	RouteID  string `gorm:"column:route_id"`       //Route id
	Order    int    `gorm:"column:order"`          //Order sort
	Config   string `gorm:"column:config"`         //Routing configuration
	IsEnable int    `gorm:"column:is_enable"`      //Plugin enable  1:on  2:off
	ModelTime
}

// TableName sets the insert table name for this struct type
func (r *RoutePlugins) TableName() string {
	return "oak_route_plugins"
}
