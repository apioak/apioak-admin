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

type Routers struct {
	ID             int    `gorm:"column:id;primary_key"`  // primary key
	ResID          string `gorm:"column:res_id"`          // Router id
	ServiceResID   string `gorm:"column:service_res_id"`  // Service id
	UpstreamResID  string `gorm:"column:upstream_res_id"` // Upstream id
	RouterName     string `gorm:"column:router_name"`     // Router name
	RequestMethods string `gorm:"column:request_methods"` // Request method
	RouterPath     string `gorm:"column:router_path"`     // Routing path
	Enable         int    `gorm:"column:enable"`          // Router enable  1:on  2:off
	Release        int    `gorm:"column:release"`         // Service release status 1:unpublished  2:to be published  3:published
	ModelTime
}

// TableName sets the insert table name for this struct type
func (r *Routers) TableName() string {
	return "oak_routers"
}

var recursionTimesRouter = 1

func (m *Routers) ModelUniqueId() (string, error) {
	generateResId, generateIdErr := utils.IdGenerate(utils.IdTypeRouter)
	if generateIdErr != nil {
		return "", generateIdErr
	}

	result := packages.GetDb().
		Table(m.TableName()).
		Where("res_id = ?", generateResId).
		Select("res_id").
		First(m)

	if result.RowsAffected == 0 {
		recursionTimesRouter = 1
		return generateResId, nil
	} else {
		if recursionTimesRouter == utils.IdGenerateMaxTimes {
			recursionTimesRouter = 1
			return "", errors.New(enums.CodeMessages(enums.IdConflict))
		}

		recursionTimesRouter++
		id, err := m.ModelUniqueId()

		if err != nil {
			return "", err
		}

		return id, nil
	}
}

func (r *Routers) RouterInfosByServiceRouterPath(
	serviceResId string,
	routerPaths []string,
	filterRouterResIds []string) ([]Routers, error) {
	routersInfos := make([]Routers, 0)
	db := packages.GetDb().
		Table(r.TableName()).
		Where("service_res_id = ?", serviceResId).
		Where("router_path IN ?", routerPaths)

	if len(filterRouterResIds) != 0 {
		db = db.Where("res_id NOT IN ?", filterRouterResIds)
	}

	err := db.Find(&routersInfos).Error

	return routersInfos, err
}

func (r *Routers) RouterDetailByResId(routerResId string) (routerDetail Routers, err error) {
	err = packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", routerResId).
		First(&routerDetail).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (r *Routers) RouterDetailByResIdServiceResId(routerResId string, serviceResId string) (routerDetail Routers) {
	db := packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", routerResId)

	if len(serviceResId) != 0 {
		db = db.Where("service_res_id = ?", serviceResId)
	}

	db.First(&routerDetail)

	return
}

func (r *Routers) RouterListByRouterResIds(routerResIds []string) ([]Routers, error) {

	routerList := make([]Routers, 0)

	err := packages.GetDb().
		Table(r.TableName()).
		Where("res_id in ?", routerResIds).
		Find(&routerList).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return routerList, nil
	}

	return routerList, err
}

func (r *Routers) RouterInfosByServiceRouterId(serviceResId string, resId string) (router Routers, err error) {
	err = packages.GetDb().
		Table(r.TableName()).
		Where("service_res_id = ?", serviceResId).
		Where("res_id = ?", resId).
		First(&router).Error

	return
}

func (r *Routers) RouterInfosByServiceIdReleaseStatus(serviceId string, releaseStatus []int) []Routers {
	routerInfos := make([]Routers, 0)
	if len(serviceId) == 0 {
		return routerInfos
	}

	db := packages.GetDb().
		Table(r.TableName()).
		Where("service_id = ?", serviceId)

	if len(releaseStatus) != 0 {
		db = db.Where("release IN ?", releaseStatus)
	}
	db.Find(&routerInfos)

	return routerInfos
}

func (r *Routers) RouterAdd(routerData Routers, upstreamData Upstreams, upstreamNodes []UpstreamNodes) (string, error) {
	routerResId, routerIdUniqueErr := r.ModelUniqueId()
	if routerIdUniqueErr != nil {
		return routerResId, routerIdUniqueErr
	}

	tx := packages.GetDb().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return routerResId, err
	}

	if len(upstreamNodes) > 0 {
		upstreamResId, err := upstreamData.ModelUniqueId()
		if err != nil {
			return routerResId, err
		}

		upstreamData.ResID = upstreamResId
		upstreamData.Name = upstreamResId
		if upstreamData.Algorithm == 0 {
			upstreamData.Algorithm = utils.LoadBalanceRoundRobin
		}
		upstreamErr := tx.Create(&upstreamData).Error

		if upstreamErr != nil {
			tx.Rollback()
			return routerResId, upstreamErr
		}

		for _, upstreamNode := range upstreamNodes {
			upstreamNodeResId, nodeErr := upstreamNode.ModelUniqueId()
			if nodeErr != nil {
				return routerResId, nodeErr
			}

			upstreamNode.ResID = upstreamNodeResId
			upstreamNode.UpstreamResID = upstreamResId

			upstreamNodeErr := tx.Create(&upstreamNode).Error
			if upstreamNodeErr != nil {
				tx.Rollback()
				return routerResId, upstreamNodeErr
			}
		}

		routerData.UpstreamResID = upstreamResId
	}

	routerData.ResID = routerResId
	if len(routerData.RouterName) == 0 {
		routerData.RouterName = routerResId
	}

	addErr := tx.Create(&routerData).Error

	if addErr != nil {
		return routerResId, addErr
	}

	return routerResId, tx.Commit().Error
}

func (r *Routers) RouterUpdate(resId string, routerData map[string]interface{}) (err error) {
	err = packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", resId).
		Updates(&routerData).Error

	return
}


func (r *Routers) RouterDelete(id string) error {

	tx := packages.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	deleteRouterError := tx.
		Table(r.TableName()).
		Where("res_id = ?", id).
		Delete(Routers{}).Error

	if deleteRouterError != nil {
		tx.Rollback()
		return deleteRouterError
	}

	return tx.Commit().Error
}

func (r *Routers) RouterListPage(serviceResId string, param *validators.ValidatorRouterList) (list []Routers, total int, listError error) {
	routersModel := Routers{}
	tx := packages.GetDb().
		Table(routersModel.TableName())

	if len(serviceResId) > 0 {
		tx = tx.Where("service_res_id = ?", serviceResId)
	}

	param.Search = strings.TrimSpace(param.Search)
	if len(param.Search) != 0 {
		search := "%" + param.Search + "%"
		orWhere := packages.GetDb().
			Or("res_id LIKE ?", search).
			Or("router_name LIKE ?", search).
			Or("router_path LIKE ?", search).
			Or("request_methods LIKE ?", strings.ToUpper(search))

		tx = tx.Where(orWhere)
	}

	if param.Enable != 0 {
		tx = tx.Where("enable = ?", param.Enable)
	}
	if param.Release != 0 {
		tx = tx.Where("`release` = ?", param.Release)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.
		Order("created_at desc")

	listError = ListPaginate(tx, &list, &param.BaseListPage)

	if len(list) == 0 {
		return
	}

	return
}

func (r *Routers) RouterUpdateName(resId string, name string) error {
	resId = strings.TrimSpace(resId)
	name = strings.TrimSpace(name)
	if (len(resId) == 0) || (len(name) == 0) {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", resId).
		Update("router_name", name).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (r *Routers) RouterSwitchEnable(id string, enable int) error {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	routerInfo, routerInfoErr := r.RouterDetailByResId(id)
	if routerInfoErr != nil {
		return routerInfoErr
	}

	releaseStatus := routerInfo.Release
	if routerInfo.Release == utils.ReleaseStatusY {
		releaseStatus = utils.ReleaseStatusT
	}

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", id).
		Updates(Routers{
			Enable:  enable,
			Release: releaseStatus}).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (r *Routers) RouterSwitchRelease(resId string, releaseStatus int) error {
	resId = strings.TrimSpace(resId)
	if len(resId) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", resId).
		Update("release", releaseStatus).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}
