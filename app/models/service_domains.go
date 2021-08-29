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

func (s *ServiceDomains) ServiceDomainIdUnique(sDomainIds map[string]string) (string, error) {
	if s.ID == "" {
		tmpID, err := utils.IdGenerate(utils.IdTypeServiceDomain)
		if err != nil {
			return "", err
		}
		s.ID = tmpID
	}

	result := packages.GetDb().Table(s.TableName()).Select("id").First(&s)
	mapId := sDomainIds[s.ID]
	if (result.RowsAffected == 0) && (s.ID != mapId) {
		sDomainId = s.ID
		sDomainIds[s.ID] = s.ID
		return sDomainId, nil
	} else {
		svcDomainId, svcErr := utils.IdGenerate(utils.IdTypeServiceDomain)
		if svcErr != nil {
			return "", svcErr
		}
		s.ID = svcDomainId
		_, err := s.ServiceDomainIdUnique(sDomainIds)
		if err != nil {
			return "", err
		}
	}

	return sDomainId, nil
}

func (s *ServiceDomains) DomainInfoByDomain(domains []string, filterServiceIds []string) []ServiceDomains {
	domainInfos := []ServiceDomains{}
	if len(filterServiceIds) == 0 {
		packages.GetDb().Table(s.TableName()).Where("domain IN ?", domains).Find(&domainInfos)
	} else {
		packages.GetDb().Table(s.TableName()).Where("domain IN ?", domains).Where("service_id NOT IN ?", filterServiceIds).Find(&domainInfos)
	}

	return domainInfos
}

func (s *ServiceDomains) DomainInfosByServiceIds(serviceIds []string) []ServiceDomains {
	domainInfos := []ServiceDomains{}
	packages.GetDb().Table(s.TableName()).Where("service_id IN ?", serviceIds).Find(&domainInfos)

	return domainInfos
}
