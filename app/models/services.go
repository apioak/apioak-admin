package models

import "time"

type Services struct {
	ID          string    `gorm:"column:id;primary_key"` //Service id
	Name        string    `gorm:"column:name"`           //Service name
	Protocol    int       `gorm:"column:protocol"`       //Protocol  0:HTTP  1:HTTPS  2:HTTP&HTTPS
	HealthCheck int       `gorm:"column:health_check"`   //Health check switch  0:off  1:on
	WebSocket   int       `gorm:"column:web_socket"`     //WebSocket  0:off  1:on
	IsEnable    int       `gorm:"column:is_enable"`      //Service enable  0:off  1:on
	LoadBalance int       `gorm:"column:load_balance"`   //Load balancing algorithm
	Timeouts    string    `gorm:"column:timeouts"`       //Time out
	CreatedAt   time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt   time.Time `gorm:"column:updated_at"`     //Update time
}

// TableName sets the insert table name for this struct type
func (s *Services) TableName() string {
	return "oak_services"
}
