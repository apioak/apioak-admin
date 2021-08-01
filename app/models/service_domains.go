package models

import "time"

type ServiceDomains struct {
	ID        string    `gorm:"column:id;primary_key"` //Domain id
	ServiceID string    `gorm:"column:service_id"`     //Service id
	Domain    string    `gorm:"column:domain"`         //Domain name
	CreatedAt time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt time.Time `gorm:"column:updated_at"`     //Update time
}

// TableName sets the insert table name for this struct type
func (s *ServiceDomains) TableName() string {
	return "oak_service_domains"
}
