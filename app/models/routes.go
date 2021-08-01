package models

import "time"

type Routes struct {
	ID             string    `gorm:"column:id;primary_key"`  //Route id
	ServiceID      string    `gorm:"column:service_id"`      //Service id
	RouteName      string    `gorm:"column:route_name"`      //Route name
	RequestMethods string    `gorm:"column:request_methods"` //Request method
	RoutePath      string    `gorm:"column:route_path"`      //Routing path
	IsEnable       int       `gorm:"column:is_enable"`       //Routing enable  0:off  1:on
	CreatedAt      time.Time `gorm:"column:created_at"`      //Creation time
	UpdatedAt      time.Time `gorm:"column:updated_at"`      //Update time
}

// TableName sets the insert table name for this struct type
func (r *Routes) TableName() string {
	return "oak_routes"
}
