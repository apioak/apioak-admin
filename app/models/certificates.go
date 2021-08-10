package models

import "time"

type Certificates struct {
	ID          string    `gorm:"column:id;primary_key"` //Certificate id
	Sni         string    `gorm:"column:sni"`            //SNI
	Certificate string    `gorm:"column:certificate"`    //Certificate content
	PrivateKey  string    `gorm:"column:private_key"`    //Private key content
	IsEnable    int       `gorm:"column:is_enable"`      //Certificate enable  1:on  2:off
	ExpiredAt   time.Time `gorm:"column:expired_at"`     //Expiration time
	ModelTime
}

// TableName sets the insert table name for this struct type
func (c *Certificates) TableName() string {
	return "oak_certificates"
}
