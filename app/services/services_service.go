package services

import (
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"fmt"
	"strings"
)

func ServiceCreate(serviceData *validators.ServiceAdd, serviceDomains *[]validators.ServiceDomainAdd, serviceNodes *[]validators.ServiceNodeAdd) error {

	serviceModel := &models.Services{}
	serviceDomainInfos := []models.ServiceDomains{}
	serviceNodeInfos := []models.ServiceNodes{}

	createServiceData := models.Services{
		Protocol:    serviceData.Protocol,
		HealthCheck: serviceData.HealthCheck,
		WebSocket:   serviceData.WebSocket,
		IsEnable:    serviceData.IsEnable,
		LoadBalance: serviceData.LoadBalance,
		Timeouts:    serviceData.Timeouts,
	}

	for _, domainInfo := range *serviceDomains {
		domain := models.ServiceDomains{
			Domain: domainInfo.Domain,
		}
		serviceDomainInfos = append(serviceDomainInfos, domain)
	}

	for _, nodeInfo := range *serviceNodes {
		ipType, err := utils.DiscernIP(nodeInfo.NodeIp)
		if err != nil {
			return err
		}
		ipTypeMap := models.IPTypeMap()
		nodeIPInfo := models.ServiceNodes{
			NodeIP:     nodeInfo.NodeIp,
			IPType:     ipTypeMap[ipType],
			NodePort:   nodeInfo.NodePort,
			NodeWeight: nodeInfo.NodeWeight,
		}
		serviceNodeInfos = append(serviceNodeInfos, nodeIPInfo)
	}

	createErr := serviceModel.ServiceAdd(&createServiceData, &serviceDomainInfos, &serviceNodeInfos)

	return createErr
}

func CheckExistDomain(domains string) error {
	serviceDomainInfo := models.ServiceDomains{}
	domainInfos := strings.Split(strings.TrimSpace(domains), ",")

	serviceDomainInfo.DomainInfoByDomain(domainInfos)
	serviceDomains := serviceDomainInfo.DomainInfoByDomain(domainInfos)

	if len(serviceDomains) == 0 {
		return nil
	}

	var existDomains = []string{}
	for _, serviceDomain := range serviceDomains {
		if len(serviceDomain.Domain) == 0 {
			continue
		}

		if len(existDomains) == 0 {
			existDomains = append(existDomains, serviceDomain.Domain)
			continue
		}

		var exist = false
		for _, existDomain := range existDomains {
			if existDomain == serviceDomain.Domain {
				exist = true
			}
		}
		if exist {
			continue
		} else {
			existDomains = append(existDomains, serviceDomain.Domain)
		}
	}

	if len(existDomains) != 0 {
		return fmt.Errorf("[" + strings.Join(existDomains, ",") + "]域名已存在")
	}

	return nil
}
