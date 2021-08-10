package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
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

var ServiceId = ""

func (service *Services)ServiceIdUnique() (string, error) {
	if ServiceId != "" {
		return ServiceId, nil
	} else {
		if service.ID == "" {
			tmpID, err := utils.IdGenerate(utils.IdTypeService)
			if err != nil {
				return "", err
			}
			service.ID = tmpID
		}

		result := packages.GetDb().First(&service);
		if result.RowsAffected == 0 {
			ServiceId = service.ID
			return ServiceId, nil
		} else {
			svcId, svcErr := utils.IdGenerate(utils.IdTypeService)
			if svcErr != nil {
				return "", svcErr
			}
			service.ID = svcId
			_, err := service.ServiceIdUnique()
			if err != nil {
				return "", err
			}
		}
	}
	return ServiceId, nil
}


func (s *Services) ServiceAdd(serviceInfo *validators.ServiceAdd, serviceDomains *[]validators.ServiceDomainAdd, serviceNodes *[]validators.ServiceNodeAdd) error {

	var serviceDomainInfos = make([]ServiceDomains, 0)
	var serviceNodeInfos = make([]ServiceNodes, 0)
	serviceId, _ := utils.IdGenerate(utils.IdTypeService)

	createServiceData := Services{
		ID:          serviceId,
		Name:        serviceId,
		Protocol:    serviceInfo.Protocol,
		HealthCheck: serviceInfo.HealthCheck,
		WebSocket:   serviceInfo.WebSocket,
		IsEnable:    serviceInfo.IsEnable,
		LoadBalance: serviceInfo.LoadBalance,
		Timeouts:    serviceInfo.Timeouts,
	}

	for _, domainInfo := range *serviceDomains {
		domainId, _ := utils.IdGenerate(utils.IdTypeMain)
		domain := ServiceDomains{
			ID:        domainId,
			ServiceID: serviceId,
			Domain:    domainInfo.Domain,
		}
		serviceDomainInfos = append(serviceDomainInfos, domain)
	}

	for _, nodeInfo := range *serviceNodes {
		nodeId, _ := utils.IdGenerate(utils.IdTypeNode)
		ipType, err := utils.DiscernIP(nodeInfo.NodeIp)
		if err != nil {
			return err
		}
		ipTypeMap := IPTypeMap()
		nodeIPInfo := ServiceNodes{
			ID:         nodeId,
			ServiceID:  serviceId,
			NodeIP:     nodeInfo.NodeIp,
			IPType:     ipTypeMap[ipType],
			NodePort:   nodeInfo.NodePort,
			NodeWeight: nodeInfo.NodeWeight,
		}
		serviceNodeInfos = append(serviceNodeInfos, nodeIPInfo)
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

	err := tx.Debug().Create(&createServiceData).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	domainErr := tx.Debug().Create(&serviceDomainInfos).Error
	if domainErr != nil {
		tx.Rollback()
		return domainErr
	}

	nodeErr := tx.Debug().Create(&serviceNodeInfos).Error
	if nodeErr != nil {
		tx.Rollback()
		return nodeErr
	}

	return tx.Commit().Error
}
