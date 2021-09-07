package models

import "time"

type Plugins struct {
	ID          string    `gorm:"column:id;primary_key"` //Plugin id
	Name        string    `gorm:"column:name"`           //Plugin name
	Tag         string    `gorm:"column:tag"`            //Plugin tag
	Icon        string    `gorm:"column:icon"`           //Plugin icon
	Description string    `gorm:"column:description"`    //Plugin description
	CreatedAt   time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt   time.Time `gorm:"column:updated_at"`     //Update time
}

// TableName sets the insert table name for this struct type
func (p *Plugins) TableName() string {
	return "oak_plugins"
}
