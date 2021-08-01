package models

import "time"

type Certificates struct {
	ID          string    `gorm:"column:id;primary_key"` //Certificate id
	Sni         string    `gorm:"column:sni"`            //SNI
	Certificate string    `gorm:"column:certificate"`    //Certificate content
	PrivateKey  string    `gorm:"column:private_key"`    //Private key content
	IsEnable    int       `gorm:"column:is_enable"`      //Certificate enable  0:off  1:on
	ExpiredAt   time.Time `gorm:"column:expired_at"`     //Expiration time
	CreatedAt   time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt   time.Time `gorm:"column:updated_at"`     //Update time
}

// TableName sets the insert table name for this struct type
func (c *Certificates) TableName() string {
	return "oak_certificates"
}
