package models

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"strings"
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

	result := packages.GetDb().
		Table(s.TableName()).
		Select("id").
		First(&s)

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

func (s *ServiceDomains) DomainInfosByDomain(domains []string, filterServiceIds []string) ([]ServiceDomains, error) {
	domainInfos := make([]ServiceDomains, 0)
	db := packages.GetDb().
		Table(s.TableName()).
		Where("domain IN ?", domains)

	if len(filterServiceIds) != 0 {
		db = db.Where("service_id NOT IN ?", filterServiceIds)
	}

	err := db.Find(&domainInfos).Error

	return domainInfos, err
}

func (s *ServiceDomains) DomainInfosByServiceIds(serviceIds []string) ([]ServiceDomains, error) {
	domainInfos := make([]ServiceDomains, 0)
	if len(serviceIds) == 0 {
		return domainInfos, nil
	}
	err := packages.GetDb().
		Table(s.TableName()).
		Where("service_id IN ?", serviceIds).
		Find(&domainInfos).Error

	if err != nil {
		return nil, err
	}

	return domainInfos, nil
}

func (s *ServiceDomains) ServiceDomainInfosLikeDomain(domain string) ([]ServiceDomains, error) {
	domainInfos := make([]ServiceDomains, 0)
	domain = strings.TrimSpace(domain)
	if len(domain) == 0 {
		return domainInfos, nil
	}

	domain = "%" + domain + "%"
	err := packages.GetDb().
		Table(s.TableName()).
		Where("domain LIKE ?", domain).
		Find(&domainInfos).Error
	
	if err != nil {
		return nil, err
	}

	return domainInfos, nil
}

func (s *ServiceDomains) DomainListByLikeSni(sni string) []ServiceDomains {
	domainList := make([]ServiceDomains, 0)
	packages.GetDb().
		Table(s.TableName()).
		Where("domain LIKE ?", "%" + sni).
		Find(&domainList)

	return domainList
}
