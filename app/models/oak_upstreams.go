package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
	"gorm.io/gorm"
)

type Upstreams struct {
	ID             int    `gorm:"column:id;primary_key"`  // primary key
	ResID          string `gorm:"column:res_id"`          // Upstream id
	Name           string `gorm:"column:name"`            // Upstream name
	Algorithm      int    `gorm:"column:algorithm"`       // Load balancing algorithm  1:round robin  2:chash
	ConnectTimeout int    `gorm:"column:connect_timeout"` // Connect timeout
	WriteTimeout   int    `gorm:"column:write_timeout"`   // Write timeout
	ReadTimeout    int    `gorm:"column:read_timeout"`    // Read timeout
	ModelTime
}

// TableName sets the insert table name for this struct type
func (u *Upstreams) TableName() string {
	return "oak_upstreams"
}

func (m *Upstreams) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeUpstream)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesServices = 1
		return generateId, nil
	} else {
		if recursionTimesServices == utils.IdGenerateMaxTimes {
			recursionTimesServices = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesServices++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (m *Upstreams) UpstreamListByResIds(upstreamResIds []string) ([]Upstreams, error) {
	upstreamList := make([]Upstreams, 0)

	err := packages.GetDb().
		Table(m.TableName()).
		Where("res_id in ?", upstreamResIds).
		Find(&upstreamList).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return upstreamList, nil
	}

	return upstreamList, err
}

func (m *Upstreams) UpstreamDetailByResId(resIds string) (detail Upstreams, err error) {
	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", resIds).
		First(&detail).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}
