package models

import "time"

type RoutePlugins struct {
	ID        string    `gorm:"column:id;primary_key"` //Plugin id
	RouteID   string    `gorm:"column:route_id"`       //Route id
	Order     int       `gorm:"column:order"`          //Order sort
	Config    string    `gorm:"column:config"`         //Routing configuration
	IsEnable  int       `gorm:"column:is_enable"`      //Plugin enable  0:off  1:on
	CreatedAt time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt time.Time `gorm:"column:updated_at"`     //Update time
}

// TableName sets the insert table name for this struct type
func (r *RoutePlugins) TableName() string {
	return "oak_route_plugins"
}
