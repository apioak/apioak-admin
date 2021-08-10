package models

type Users struct {
	ID       string `gorm:"column:id;primary_key"` //User iD
	Name     string `gorm:"column:name"`           //User name
	Password string `gorm:"column:password"`       //Password
	Email    string `gorm:"column:email"`          //Email
	ModelTime
}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "oak_users"
}
