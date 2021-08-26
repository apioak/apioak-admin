package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
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
}

// TableName sets the insert table name for this struct type
func (s *Services) TableName() string {
	return "oak_services"
}

var sId = ""

func (service *Services) ServiceIdUnique(sIds map[string]string) (string, error) {
	if service.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeService)
		if err != nil {
			return "", err
		}
		service.ID = tmpID
	}

	result := packages.GetDb().Table(service.TableName()).Select("id").First(&service)
	mapId := sIds[service.ID]
	if (result.RowsAffected == 0) && (service.ID != mapId) {
		sId = service.ID
		sIds[service.ID] = service.ID
		return sId, nil
	} else {
		svcId, svcErr := utils.IdGenerate(utils.IdTypeService)
		if svcErr != nil {
			return "", svcErr
		}
		service.ID = svcId
		_, err := service.ServiceIdUnique(sIds)
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
