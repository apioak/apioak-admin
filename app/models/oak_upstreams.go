package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type Upstreams struct {
	ID             int    `gorm:"column:id;primary_key"`  // primary key
	ResID          string `gorm:"column:res_id"`          // Upstream id
	Name           string `gorm:"column:name"`            // Upstream name
	Algorithm      int    `gorm:"column:algorithm"`       // Load balancing algorithm  1:round robin  2:chash
	ConnectTimeout int    `gorm:"column:connect_timeout"` // Connect timeout
	WriteTimeout   int    `gorm:"column:write_timeout"`   // Write timeout
	ReadTimeout    int    `gorm:"column:read_timeout"`    // Read timeout
	Enable         int    `gorm:"column:enable"`          // Enable  1:on  2:off
	Release        int    `gorm:"column:release"`         // Release status 1:unpublished  2:to be published  3:published
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

func (m Upstreams) UpstreamListPage(resIds []string, request *validators.UpstreamList) (list []Upstreams, total int, err error) {
	tx := packages.GetDb().Table(m.TableName())

	if len(resIds) != 0 {
		tx.Where("res_id IN ?", resIds)
	}

	if request.Search != "" {
		search := "%" + request.Search + "%"
		orWhere := packages.GetDb().
			Or("res_id LIKE ?", search).
			Or("name LIKE ?", search)

		tx = tx.Where(orWhere)
	}

	if request.Algorithm != 0 {
		tx.Where("algorithm = ?", request.Algorithm)
	}
	if request.Enable != 0 {
		tx = tx.Where("enable = ?", request.Enable)
	}
	if request.Release != 0 {
		tx = tx.Where("`release` = ?", request.Release)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		err = countError
		return
	}

	tx = tx.
		Order("created_at desc")

	err = ListPaginate(tx, &list, &request.BaseListPage)

	if len(list) == 0 {
		return
	}

	 return
}

func (m Upstreams) UpstreamInfosByNames (names []string, filterResIds []string) (list []Upstreams, err error) {
	list = make([]Upstreams, 0)

	if len(names) == 0 {
		return
	}

	db := packages.GetDb().Table(m.TableName()).
		Where("name IN ?", names)

	if len(filterResIds) != 0 {
		db = db.Where("res_id NOT IN ?", filterResIds)
	}

	err = db.Find(&list).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (m Upstreams) UpstreamAdd(upstreamData Upstreams, upstreamNodes []UpstreamNodes) (resId string, err error) {
	resId, err = m.ModelUniqueId()
	if err != nil {
		return
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {
		upstreamData.ResID = resId
		if upstreamData.Name == "" {
			upstreamData.Name = resId
		}

		if upstreamData.Algorithm == 0 {
			upstreamData.Algorithm = utils.LoadBalanceRoundRobin
		}

		err = tx.Create(&upstreamData).Error
		if err != nil {
			return err
		}

		if len(upstreamNodes) == 0 {
			return nil
		}

		for key, upstreamNodeInfo := range upstreamNodes {
			nodeResId, nodeResIdErr := upstreamNodeInfo.ModelUniqueId()
			if nodeResIdErr != nil {
				return nodeResIdErr
			}

			upstreamNodes[key].UpstreamResID = resId
			upstreamNodes[key].ResID = nodeResId
		}

		err = tx.Create(&upstreamNodes).Error
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (m Upstreams) UpstreamReleaseNameList() (list []ResIdNameItem, err error) {
	list = make([]ResIdNameItem, 0)

	err = packages.GetDb().Table(m.TableName()).
		Where("enable = ?", utils.EnableOn).
		Where("`release` = ?", utils.ReleaseStatusY).
		Find(&list).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (m Upstreams) UpstreamUpdateName(resId string, name string) (err error) {
	name = strings.TrimSpace(name)

	if (len(resId) == 0) || (len(name) == 0) {
		return errors.New(enums.CodeMessages(enums.ParamsError))
	}

	err = packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", resId).
		Update("name", name).Error

	return
}


