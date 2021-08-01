package models

import "time"

type ServiceNodes struct {
	ID         string    `gorm:"column:id;primary_key"` //Service node id
	ServiceID  string    `gorm:"column:service_id"`     //Service id
	NodeIP     string    `gorm:"column:node_ip"`        //Node IP
	IPType     int       `gorm:"column:ip_type"`        //IP Type  0:IPV4  1:IPV6
	NodePort   int       `gorm:"column:node_port"`      //Node port
	NodeWeight int       `gorm:"column:node_weight"`    //Node weight
	CreatedAt  time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt  time.Time `gorm:"column:updated_at"`     //Update time
}

// TableName sets the insert table name for this struct type
func (s *ServiceNodes) TableName() string {
	return "oak_service_nodes"
}
