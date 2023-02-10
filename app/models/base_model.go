package models

import (
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"gorm.io/gorm"
	"time"
)

type ModelTime struct {
	CreatedAt time.Time `gorm:"column:created_at"` // Creation time
	UpdatedAt time.Time `gorm:"column:updated_at"` // Update time
}

type ResIdNameItem struct {
	ResId string `json:"res_id"`
	Name  string `json:"name"`
}

func ListCount(db *gorm.DB, total *int) error {
	var count int64
	countError := db.Count(&count).Error
	if countError != nil {
		return countError
	}
	*total = int(count)

	return nil
}

func ListPaginate(db *gorm.DB, list interface{}, page *validators.BaseListPage) error {
	if page.Page <= utils.Page {
		page.Page = utils.Page
	}

	switch {
	case page.PageSize > utils.MaxPageSize:
		page.PageSize = utils.MaxPageSize
	case page.PageSize <= 0:
		page.PageSize = utils.PageSize
	}

	offset := (page.Page - 1) * page.PageSize
	listError := db.Limit(page.PageSize).Offset(offset).Find(list).Error
	if listError != nil {
		return listError
	}

	return nil
}
