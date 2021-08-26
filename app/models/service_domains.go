package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
)

type ServiceDomains struct {
	ID        string `gorm:"column:id;primary_key"` //Domain id
	ServiceID string `gorm:"column:service_id"`     //Service id
	Domain    string `gorm:"column:domain"`         //Domain name
	ModelTime
}

// TableName sets the insert table name for this struct type
func (s *ServiceDomains) TableName() string {
	return "oak_service_domains"
}

var sDomainId = ""

func (serviceDomain *ServiceDomains) ServiceDomainIdUnique(sDomainIds map[string]string) (string, error) {
	if serviceDomain.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeServiceDomain)
		if err != nil {
			return "", err
		}
		serviceDomain.ID = tmpID
	}

	result := packages.GetDb().Table(serviceDomain.TableName()).Select("id").First(&serviceDomain)
	mapId := sDomainIds[serviceDomain.ID]
	if (result.RowsAffected == 0) && (serviceDomain.ID != mapId) {
		sDomainId = serviceDomain.ID
		sDomainIds[serviceDomain.ID] = serviceDomain.ID
		return sDomainId, nil
	} else {
		svcDomainId, svcErr := utils.IdGenerate(utils.IdTypeServiceDomain)
		if svcErr != nil {
			return "", svcErr
		}
		serviceDomain.ID = svcDomainId
		_, err := serviceDomain.ServiceDomainIdUnique(sDomainIds)
		if err != nil {
			return "", err
		}
	}

	return sDomainId, nil
}

func (serviceDomain *ServiceDomains) DomainInfoByDomain(domains []string) []ServiceDomains {
	serviceDomains := &ServiceDomains{}

	domainInfos := []ServiceDomains{}
	packages.GetDb().Table(serviceDomains.TableName()).Where("domain IN ?", domains).Find(&domainInfos)

	return domainInfos
}
