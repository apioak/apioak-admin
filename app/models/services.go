package models

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"strings"
)

type Services struct {
	ID          string `gorm:"column:id;primary_key"` //Service id
	Name        string `gorm:"column:name"`           //Service name
	Protocol    int    `gorm:"column:protocol"`       //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	HealthCheck int    `gorm:"column:health_check"`   //Health check switch  1:on  2:off
	WebSocket   int    `gorm:"column:web_socket"`     //WebSocket  1:on  2:off
	IsEnable    int    `gorm:"column:is_enable"`      //Service enable  1:on  2:off
	LoadBalance int    `gorm:"column:load_balance"`   //Load balancing algorithm
	Timeouts    string `gorm:"column:timeouts"`       //Time out
	ModelTime
	Domains []ServiceDomains `gorm:"foreignKey:ServiceID"`
	Nodes   []ServiceNodes   `gorm:"foreignKey:ServiceID"`
}

// TableName sets the insert table name for this struct type
func (s *Services) TableName() string {
	return "oak_services"
}

var sId = ""

func (s *Services) ServiceIdUnique(sIds map[string]string) (string, error) {
	if s.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeService)
		if err != nil {
			return "", err
		}
		s.ID = tmpID
	}

	result := packages.GetDb().Table(s.TableName()).Select("id").First(&s)
	mapId := sIds[s.ID]
	if (result.RowsAffected == 0) && (s.ID != mapId) {
		sId = s.ID
		sIds[s.ID] = s.ID
		return sId, nil
	} else {
		svcId, svcErr := utils.IdGenerate(utils.IdTypeService)
		if svcErr != nil {
			return "", svcErr
		}
		s.ID = svcId
		_, err := s.ServiceIdUnique(sIds)
		if err != nil {
			return "", err
		}
	}

	return sId, nil
}

func (s *Services) ServiceAdd(serviceInfo *Services, serviceDomains *[]ServiceDomains, serviceNodes *[]ServiceNodes) error {

	tpmIds := map[string]string{}
	serviceId, serviceIdUniqueErr := s.ServiceIdUnique(tpmIds)
	if serviceIdUniqueErr != nil {
		return serviceIdUniqueErr
	}

	tx := packages.GetDb().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	serviceInfo.ID = serviceId
	serviceInfo.Name = serviceId
	createServiceErr := tx.Create(&serviceInfo).Error
	if createServiceErr != nil {
		tx.Rollback()
		return createServiceErr
	}

	for _, serviceDomain := range *serviceDomains {
		domainId, domainIdErr := serviceDomain.ServiceDomainIdUnique(tpmIds)
		if domainIdErr != nil {
			return domainIdErr
		}

		serviceDomain.ID = domainId
		serviceDomain.ServiceID = serviceId
		domainErr := tx.Create(&serviceDomain).Error
		if domainErr != nil {
			tx.Rollback()
			return domainErr
		}
	}

	for _, serviceNode := range *serviceNodes {
		nodeId, nodeIdErr := serviceNode.ServiceNodeIdUnique(tpmIds)
		if nodeIdErr != nil {
			return nodeIdErr
		}

		serviceNode.ID = nodeId
		serviceNode.ServiceID = serviceId
		nodeErr := tx.Create(&serviceNode).Error
		if nodeErr != nil {
			tx.Rollback()
			return nodeErr
		}
	}

	return tx.Commit().Error
}

func (s *Services) ServiceInfoById(id string) Services {
	serviceInfo := Services{}
	packages.GetDb().Table(s.TableName()).Where("id = ?", id).First(&serviceInfo)
	return serviceInfo
}

func (s *Services) ServiceUpdate(
	id string,
	serviceInfo *Services,
	serviceDomains *[]ServiceDomains,
	serviceNodes *[]ServiceNodes,
	updateNodes *[]ServiceNodes,
	deleteDomainIds []string,
	deleteNodeIds []string) error {

	tx := packages.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	updateError := tx.Table(s.TableName()).Where("id = ?", id).Updates(serviceInfo).Error
	if updateError != nil {
		tx.Rollback()
		return updateError
	}

	if len(deleteDomainIds) != 0 {
		serviceDomainsModel := ServiceDomains{}
		domainDeleteError := tx.Table(serviceDomainsModel.TableName()).Where("id IN ?", deleteDomainIds).Delete(&serviceDomainsModel).Error
		if domainDeleteError != nil {
			tx.Rollback()
			return domainDeleteError
		}
	}

	if len(deleteNodeIds) != 0 {
		serviceNodesModel := ServiceNodes{}
		nodeDeleteError := tx.Table(serviceNodesModel.TableName()).Where("id IN ?", deleteNodeIds).Delete(&serviceNodesModel).Error
		if nodeDeleteError != nil {
			tx.Rollback()
			return nodeDeleteError
		}
	}

	if len(*serviceDomains) > 0 {
		tpmIds := map[string]string{}
		for _, serviceDomain := range *serviceDomains {
			domainId, domainIdErr := serviceDomain.ServiceDomainIdUnique(tpmIds)
			if domainIdErr != nil {
				return domainIdErr
			}

			serviceDomain.ID = domainId
			serviceDomain.ServiceID = id
			domainErr := tx.Create(&serviceDomain).Error
			if domainErr != nil {
				tx.Rollback()
				return domainErr
			}
		}
	}

	if len(*updateNodes) > 0 {
		for _, updateNode := range *updateNodes {
			updateNodeError := tx.Table(updateNode.TableName()).Where("id = ?", updateNode.ID).Updates(ServiceNodes{NodeWeight: updateNode.NodeWeight}).Error
			if updateNodeError != nil {
				tx.Rollback()
				return updateNodeError
			}
		}
	}

	if len(*serviceNodes) > 0 {
		tpmIds := map[string]string{}
		for _, serviceNode := range *serviceNodes {
			nodeId, nodeIdErr := serviceNode.ServiceNodeIdUnique(tpmIds)
			if nodeIdErr != nil {
				return nodeIdErr
			}

			serviceNode.ID = nodeId
			serviceNode.ServiceID = id
			nodeErr := tx.Create(&serviceNode).Error
			if nodeErr != nil {
				tx.Rollback()
				return nodeErr
			}
		}
	}

	return tx.Commit().Error
}

func (s *Services) ServiceDomainNodeByIds(serviceIds []string) ([]Services, error) {
	serviceInfos := make([]Services, 0)
	err := errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	if len(serviceIds) == 0 {
		return serviceInfos, err
	}

	packages.GetDb().Table(s.TableName()).
		Where("id IN ?", serviceIds).
		Preload("Domains").
		Preload("Nodes").
		Order("updated_at desc").
		Find(&serviceInfos)
	if len(serviceInfos) == 0 {
		err = errors.New(enums.CodeMessages(enums.ServiceNull))
	} else {
		err = nil
	}

	return serviceInfos, err
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

	deleteServiceError := tx.Table(s.TableName()).Where("id = ?", id).Delete(ServiceNodes{}).Error
	if deleteServiceError != nil {
		tx.Rollback()
		return deleteServiceError
	}

	serviceDomainsModel := ServiceDomains{}
	deleteDomainError := tx.Table(serviceDomainsModel.TableName()).Where("service_id = ?", id).Delete(serviceDomainsModel).Error
	if deleteDomainError != nil {
		tx.Rollback()
		return deleteDomainError
	}

	serviceNodesModel := ServiceNodes{}
	deleteNodeError := tx.Table(serviceNodesModel.TableName()).Where("service_id = ?", id).Delete(serviceNodesModel).Error
	if deleteNodeError != nil {
		tx.Rollback()
		return deleteNodeError
	}

	return tx.Commit().Error
}

func (s *Services) ServiceAllInfosListPage(
	serviceIds []string,
	param *validators.ServiceList) (list []Services, total int, listError error) {
	tx := packages.GetDb().Table(s.TableName())
	if len(serviceIds) != 0 {
		tx = tx.Where("id IN ?", serviceIds)
	}
	if param.Protocol != 0 {
		tx = tx.Where("protocol = ?", param.Protocol)
	}
	if param.IsEnable != 0 {
		tx = tx.Where("is_enable = ?", param.IsEnable)
	}
	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx = tx.Preload("Domains").Order("updated_at desc")
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
	err := packages.GetDb().Table(s.TableName()).
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

	updateErr := packages.GetDb().Table(s.TableName()).Where("id = ?", id).Update("name", name).Error
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

	updateErr := packages.GetDb().Table(s.TableName()).Where("id = ?", id).Update("is_enable", enable).Error
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (s *Services) ServiceSwitchWebsocket(id string, webSocket int) error {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().Table(s.TableName()).Where("id = ?", id).Update("web_socket", webSocket).Error
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (s *Services) ServiceSwitchHealthCheck(id string, healthCheck int) error {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	}

	updateErr := packages.GetDb().Table(s.TableName()).Where("id = ?", id).Update("health_check", healthCheck).Error
	if updateErr != nil {
		return updateErr
	}

	return nil
}
