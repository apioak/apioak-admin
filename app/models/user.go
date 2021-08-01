package models

import "time"

type users struct {
	ID        string    `gorm:"column:id;primary_key"` //User iD
	Name      string    `gorm:"column:name"`           //User name
	Password  string    `gorm:"column:password"`       //Password
	Email     string    `gorm:"column:email"`          //Email
	CreatedAt time.Time `gorm:"column:created_at"`     //Creation time
	UpdatedAt time.Time `gorm:"column:updated_at"`     //Update time
}

// TableName sets the insert table name for this struct type
func (u *users) TableName() string {
	return "oak_users"
}

