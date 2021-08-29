package services

import (
	"apioak-admin/app/models"
	"apioak-admin/app/validators"
)

func GetToOperateDomains(serviceId string, updateDomains *[]validators.ServiceDomainAddUpdate) ([]models.ServiceDomains, []string) {
	serviceDomainsModel := models.ServiceDomains{}
	serviceExistDomains, err := serviceDomainsModel.DomainInfosByServiceIds([]string{serviceId})
	if err != nil {
		// @todo 处理错误
	}

	updateDomainsMap := make(map[string]string)
	for _, updateDomain := range *updateDomains {
		updateDomainsMap[updateDomain.Domain] = updateDomain.Domain
	}

	existDomainsMap := make(map[string]string)
	for _, existDomain := range serviceExistDomains {
		existDomainsMap[existDomain.Domain] = existDomain.Domain
	}

	addDomains := []models.ServiceDomains{}
	for _, updateDomain := range *updateDomains {
		_, exist := existDomainsMap[updateDomain.Domain]
		if exist {
			continue
		}
		domain := models.ServiceDomains{
			ServiceID: serviceId,
			Domain:    updateDomain.Domain,
		}
		addDomains = append(addDomains, domain)
	}

	deleteDomainIds := []string{}
	for _, existDomain := range serviceExistDomains {
		_, exist := updateDomainsMap[existDomain.Domain]
		if !exist {
			deleteDomainIds = append(deleteDomainIds, existDomain.ID)
		}
	}

	return addDomains, deleteDomainIds
}
