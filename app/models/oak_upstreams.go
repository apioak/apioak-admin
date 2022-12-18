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
func (u Upstreams) TableName() string {
	return "oak_upstreams"
}

func (m *Upstreams) ModelUniqueId() (generateId string, err error) {
	generateId, err = utils.IdGenerate(utils.IdTypeUpstream)
	if err != nil {
		return
	}

	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateId).
		Select("res_id").
		First(m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err == nil {
		recursionTimesServices = 1
		return
	} else {
		if recursionTimesServices == utils.IdGenerateMaxTimes {
			recursionTimesServices = 1
			err = errors.New(enums.CodeMessages(enums.IdConflict))
			return
		}

		recursionTimesServices++
		generateId, err = m.ModelUniqueId()
		if err != nil {
			return
		}

		return
	}
}

func (m Upstreams) UpstreamListByResIds(upstreamResIds []string) ([]Upstreams, error) {
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

func (m Upstreams) UpstreamDetailByResId(resId string) (detail Upstreams, err error) {
	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", resId).
		First(&detail).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (m Upstreams) UpstreamUpdate(resIds string, upstreamData Upstreams) (err error) {
	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", resIds).
		Updates(&upstreamData).Error

	return
}


