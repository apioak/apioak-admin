package models

import (
	"time"
)

type ModelTime struct {
	CreatedAt time.Time `gorm:"column:created_at"` //Creation time
	UpdatedAt time.Time `gorm:"column:updated_at"` //Update time
}
