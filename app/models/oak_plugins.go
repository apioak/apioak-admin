package models

type Plugins struct {
	ID          string `gorm:"column:id;primary_key"` //Plugin id
	Name        string `gorm:"column:name"`           //Plugin name
	Tag         string `gorm:"column:tag"`            //Plugin tag
	Icon        string `gorm:"column:icon"`           //Plugin icon
	Description string `gorm:"column:description"`    //Plugin description
	ModelTime
}

// TableName sets the insert table name for this struct type
func (p *Plugins) TableName() string {
	return "oak_plugins"
}
