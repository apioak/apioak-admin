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

func (r *Routers) RouterInfoById(routerId string) (Routers, error) {
	routerInfo := Routers{}
	err := packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", routerId).
		First(&routerInfo).Error

	return routerInfo, err
}

func (r *Routers) RouterInfoByResIdServiceResId(routerResId string, serviceResId string) Routers {
	routerInfo := Routers{}
	db := packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", routerResId)

	if len(serviceResId) != 0 {
		db = db.Where("service_res_id = ?", serviceResId)
	}

	db.First(&routerInfo)

	return routerInfo
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

func (r *Routers) RouterUpdate(id string, routerData Routers) error {
	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", id).
		Updates(&routerData).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

// func (r *Routers) RouterCopy(routerData Routers, routerPlugins []RouterPlugins) (string, error) {
// 	tx := packages.GetDb().Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()
//
// 	if err := tx.Error; err != nil {
// 		return "", err
// 	}
//
// 	routerId, routerIdUniqueErr := r.ModelUniqueId()
// 	if routerIdUniqueErr != nil {
// 		return routerId, routerIdUniqueErr
// 	}
//
// 	routerData.ResID = routerId
// 	routerData.RouterName = routerId
//
// 	addErr := tx.
// 		Table(r.TableName()).
// 		Create(&routerData).Error
//
// 	if addErr != nil {
// 		tx.Rollback()
// 		return routerId, addErr
// 	}
//
// 	if len(routerPlugins) != 0 {
// 		routerPluginModel := RouterPlugins{}
// 		for k, _ := range routerPlugins {
// 			routerPluginId, routerPluginIdErr := routerPluginModel.ModelUniqueId()
// 			if routerPluginIdErr != nil {
// 				tx.Rollback()
// 				return routerId, routerPluginIdErr
// 			}
// 			±
// 			routerPlugins[k].ID = routerPluginId
// 			routerPlugins[k].RouterID = routerId
// 			routerPlugins[k].ReleaseStatus = utils.ReleaseStatusU
// 			routerPlugins[k].CreatedAt = time.Now()
// 			routerPlugins[k].UpdatedAt = time.Now()
// 		}
//
// 		addRouterPluginErr := tx.
// 			Table(routerPluginModel.TableName()).
// 			Create(&routerPlugins).Error
//
// 		if addRouterPluginErr != nil {
// 			tx.Rollback()
// 			return routerId, addRouterPluginErr
// 		}
// 	}
//
// 	return routerId, tx.Commit().Error
// }

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

type RouterPluginListItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Icon          string `json:"icon"`
	Type          int    `json:"type"`
	Config        string `json:"config"`
	IsEnable      int    `json:"is_enable"`
	ReleaseStatus int    `json:"release_status"`
}

type RouterListItem struct {
	ResID          string                 `json:"res_id"`
	ServiceResID   string                 `json:"service_res_id"`
	ServiceName    string                 `json:"service_name"`
	RouterName     string                 `json:"router_name"`
	RequestMethods string                 `json:"request_methods"`
	RouterPath     string                 `json:"router_path"`
	Enable         int                    `json:"enable"`
	Release        int                    `json:"release"`
	Plugins        []RouterPluginListItem `json:"plugins"`
}

func (r *Routers) RouterListPage(serviceResId string, param *validators.ValidatorRouterList) (list []RouterListItem, total int, listError error) {
	routersList := make([]Routers, 0)

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
		tx = tx.Where("release = ?", param.Release)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.
		Order("created_at desc")

	listError = ListPaginate(tx, &routersList, &param.BaseListPage)

	if len(routersList) == 0 {
		return
	}

	// @todo 列表中需要补充服务名称，路由插件的数据列表

	for _, routersInfo := range routersList {
		routerPluginListItem := RouterListItem{
			ResID:          routersInfo.ResID,
			ServiceResID:   routersInfo.ServiceResID,
			RouterName:     routersInfo.RouterName,
			RequestMethods: routersInfo.RequestMethods,
			RouterPath:     routersInfo.RouterPath,
			Enable:         routersInfo.Enable,
			Release:        routersInfo.Release,
		}

		list = append(list, routerPluginListItem)
	}

	return
}

func (r *Routers) RouterUpdateName(id string, name string) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	if (len(id) == 0) || (len(name) == 0) {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().
		Table(r.TableName()).
		Where("res_id = ?", id).
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

	routerInfo, routerInfoErr := r.RouterInfoById(id)
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
