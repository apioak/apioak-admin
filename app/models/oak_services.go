package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Services struct {
	ID       int    `gorm:"column:id;primary_key"` //Service id
	ResID    string `gorm:"column:res_id"`         //ServiceResID
	Name     string `gorm:"column:name"`           //Service name
	Protocol int    `gorm:"column:protocol"`       //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	Enable   int    `gorm:"column:enable"`         //Service enable  1:on  2:off
	Release  int    `gorm:"column:release"`        //Service release status 1:unpublished  2:to be published  3:published
	ModelTime
}

type ServiceDomainNode struct {
	Services
	Domains []ServiceDomains `gorm:"foreignKey:ServiceID"` //domains
	Nodes   []ServiceNodes   `gorm:"foreignKey:ServiceID"` //nodes(upstreams)
}

// TableName sets the insert table name for this struct type
func (s *Services) TableName() string {
	return "oak_services"
}

var recursionTimesServices = 1

func (m *Services) ModelUniqueId() (string, error) {
	generateId, generateIdErr := utils.IdGenerate(utils.IdTypeService)
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

func (s *Services) ServiceAdd(serviceInfo *Services, serviceDomains []string) (string, error) {

	serviceId, err := s.ModelUniqueId()
	if err != nil {
		return serviceId, err
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {
		serviceInfo.ResID = serviceId
		if serviceInfo.Name == "" {
			serviceInfo.Name = serviceId
		}

		err := tx.Create(&serviceInfo).Error

		if err != nil {
			packages.GetLogger().Error("create services error")
			return err
		}

		serviceDomainsParam := []ServiceDomains{}
		for _, serviceDomain := range serviceDomains {
			domainId, err := (&ServiceDomains{}).ModelUniqueId()
			if err != nil {
				return err
			}

			serviceDomainsParam = append(serviceDomainsParam, ServiceDomains{
				Domain:    serviceDomain,
				ResID:     domainId,
				ServiceID: serviceId,
			})
		}

		err = tx.Create(&serviceDomainsParam).Error

		if err != nil {
			packages.GetLogger().Error("create services domain error")
			return err
		}

		return nil
	})

	fmt.Println(err)

	return serviceId, err
}

func (s *Services) ServiceUpdateColumnsById(id string, serviceInfo *Services) error {
	return packages.GetDb().
		Table(s.TableName()).
		Where("id = ?", id).
		Updates(serviceInfo).Error
}

func (s *Services) ServiceInfoById(id string) Services {
	serviceInfo := Services{}
	packages.GetDb().
		Table(s.TableName()).
		Where("id = ?", id).
		First(&serviceInfo)

	return serviceInfo
}

func (s *Services) ServiceUpdate(
	id string,
	serviceInfo *Services,
	serviceDomains *[]ServiceDomains,
) error {

	tx := packages.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	updateError := tx.
		Table(s.TableName()).
		Where("id = ?", id).
		Updates(serviceInfo).Error

	if updateError != nil {
		tx.Rollback()
		return updateError
	}

	if len(*serviceDomains) > 0 {
		for _, serviceDomain := range *serviceDomains {
			domainId, domainIdErr := serviceDomain.ModelUniqueId()
			if domainIdErr != nil {
				return domainIdErr
			}

			serviceDomain.ResID = domainId
			serviceDomain.ServiceID = id
			domainErr := tx.
				Create(&serviceDomain).Error

			if domainErr != nil {
				tx.Rollback()
				return domainErr
			}
		}
	}
	return tx.Commit().Error
}

func (s *Services) ServiceDomainNodeById(serviceId string) (ServiceDomainNode, error) {
	serviceDomainNode := ServiceDomainNode{}
	err := errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	if len(serviceId) == 0 {
		return serviceDomainNode, err
	}

	packages.GetDb().
		Table(serviceDomainNode.TableName()).
		Where("id = ?", serviceId).
		Preload("Domains").
		Preload("Nodes").
		Order("updated_at desc").
		First(&serviceDomainNode)

	if serviceDomainNode.ResID != serviceId {
		err = errors.New(enums.CodeMessages(enums.ServiceNull))
	} else {
		err = nil
	}

	return serviceDomainNode, err
}

func (s *Services) ServiceDomainNodeByIds(serviceIds []string) ([]ServiceDomainNode, error) {
	serviceDomainNode := ServiceDomainNode{}
	serviceDomainNodes := make([]ServiceDomainNode, 0)
	err := errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	if len(serviceIds) == 0 {
		return serviceDomainNodes, err
	}

	packages.GetDb().
		Table(serviceDomainNode.TableName()).
		Where("id IN ?", serviceIds).
		Preload("Domains").
		Preload("Nodes").
		Order("updated_at desc").
		Find(&serviceDomainNodes)

	if len(serviceDomainNodes) == 0 {
		err = errors.New(enums.CodeMessages(enums.ServiceNull))
	} else {
		err = nil
	}

	return serviceDomainNodes, err
}

func (s *Services) ServiceDelete(id string) error {

	tx := packages.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	deleteServiceError := tx.
		Table(s.TableName()).
		Where("id = ?", id).
		Delete(ServiceNodes{}).Error

	if deleteServiceError != nil {
		tx.Rollback()
		return deleteServiceError
	}

	serviceDomainsModel := ServiceDomains{}
	deleteDomainError := tx.
		Table(serviceDomainsModel.TableName()).
		Where("service_id = ?", id).
		Delete(serviceDomainsModel).Error

	if deleteDomainError != nil {
		tx.Rollback()
		return deleteDomainError
	}

	serviceNodesModel := ServiceNodes{}
	deleteNodeError := tx.
		Table(serviceNodesModel.TableName()).
		Where("service_id = ?", id).
		Delete(serviceNodesModel).Error

	if deleteNodeError != nil {
		tx.Rollback()
		return deleteNodeError
	}

	serviceRouteIds := make([]string, 0)
	routeModel := Routes{}
	tx.
		Table(routeModel.TableName()).
		Where("service_id = ?", id).
		Pluck("id", &serviceRouteIds)

	if len(serviceRouteIds) > 0 {
		deleteRouteError := tx.
			Table(routeModel.TableName()).
			Where("id IN ?", serviceRouteIds).
			Delete(routeModel).Error

		if deleteRouteError != nil {
			tx.Rollback()
			return deleteRouteError
		}

		routePluginModel := RoutePlugins{}
		deleteRoutePluginError := tx.
			Table(routePluginModel.TableName()).
			Where("route_id IN ?", serviceRouteIds).
			Delete(routePluginModel).Error

		if deleteRoutePluginError != nil {
			tx.Rollback()
			return deleteRoutePluginError
		}
	}

	return tx.Commit().Error
}

func (s *Services) ServiceAllInfosListPage(
	serviceIds []string,
	param *validators.ServiceList) (list []ServiceDomainNode, total int, listError error) {
	serviceDomainNode := ServiceDomainNode{}

	tx := packages.GetDb().
		Table(serviceDomainNode.TableName())

	if len(serviceIds) != 0 {
		tx = tx.Where("id IN ?", serviceIds)
	}
	if param.Protocol != 0 {
		tx = tx.Where("protocol = ?", param.Protocol)
	}
	if param.Enable != 0 {
		tx = tx.Where("is_enable = ?", param.Enable)
	}
	if param.Release != 0 {
		tx = tx.Where("release_status = ?", param.Release)
	}

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.
		Preload("Domains").
		Order("created_at desc")

	listError = ListPaginate(tx, &list, &param.BaseListPage)
	return
}

func (s *Services) ServiceInfosLikeIdName(idOrName string) ([]Services, error) {
	serviceInfos := make([]Services, 0)
	idOrName = strings.TrimSpace(idOrName)
	if len(idOrName) == 0 {
		return serviceInfos, nil
	}

	idOrName = "%" + idOrName + "%"
	err := packages.GetDb().
		Table(s.TableName()).
		Where("id LIKE ?", idOrName).
		Or("name LIKE ?", idOrName).
		Find(&serviceInfos).Error

	if err != nil {
		return nil, err
	}

	return serviceInfos, nil
}

func (s *Services) ServiceUpdateName(id string, name string) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	if (len(id) == 0) || (len(name) == 0) {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().
		Table(s.TableName()).
		Where("id = ?", id).
		Update("name", name).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (s *Services) ServiceSwitchEnable(id string, enable int) error {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	serviceInfo := s.ServiceInfoById(id)
	releaseStatus := serviceInfo.Release
	if serviceInfo.Release == utils.ReleaseStatusY {
		releaseStatus = utils.ReleaseStatusT
	}

	updateErr := packages.GetDb().
		Table(s.TableName()).
		Where("id = ?", id).
		Updates(Services{
			Enable:  enable,
			Release: releaseStatus}).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}

//func (s *Services) ServiceSwitchWebsocket(id string, webSocket int) error {
//	id = strings.TrimSpace(id)
//	if len(id) == 0 {
//		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
//	}
//
//	serviceInfo := s.ServiceInfoById(id)
//	releaseStatus := serviceInfo.Release
//	if serviceInfo.ReleaseStatus == utils.ReleaseStatusY {
//		releaseStatus = utils.ReleaseStatusT
//	}
//
//	updateErr := packages.GetDb().
//		Table(s.TableName()).
//		Where("id = ?", id).
//		Updates(Services{
//			WebSocket:     webSocket,
//			ReleaseStatus: releaseStatus}).Error
//
//	if updateErr != nil {
//		return updateErr
//	}
//
//	return nil
//}
//
//func (s *Services) ServiceSwitchHealthCheck(id string, healthCheck int) error {
//	id = strings.TrimSpace(id)
//	if len(id) == 0 {
//		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
//	}
//
//	serviceInfo := s.ServiceInfoById(id)
//	releaseStatus := serviceInfo.ReleaseStatus
//	if serviceInfo.ReleaseStatus == utils.ReleaseStatusY {
//		releaseStatus = utils.ReleaseStatusT
//	}
//
//	updateErr := packages.GetDb().
//		Table(s.TableName()).
//		Where("id = ?", id).
//		Updates(Services{
//			HealthCheck:   healthCheck,
//			ReleaseStatus: releaseStatus}).Error
//
//	if updateErr != nil {
//		return updateErr
//	}
//
//	return nil
//}

func (s *Services) ServiceSwitchRelease(id string, release int) error {
	if len(id) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().
		Table(s.TableName()).
		Where("id = ?", id).
		Update("release_status", release).Error

	if updateErr != nil {
		return updateErr
	}

	return nil
}
