package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"fmt"
	"strings"
)

func CheckExistDomain(domains string, filterServiceIds []string) error {
	serviceDomainInfo := models.ServiceDomains{}
	domainInfos := strings.Split(strings.TrimSpace(domains), ",")
	serviceDomains := serviceDomainInfo.DomainInfoByDomain(domainInfos, filterServiceIds)

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
		if !exist {
			existDomains = append(existDomains, serviceDomain.Domain)
		}
	}

	if len(existDomains) != 0 {
		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.ServiceDomainExist), strings.Join(existDomains, ",")))
	}

	return nil
}

func ServiceCreate(
	serviceData *validators.ServiceAddUpdate,
	serviceDomains *[]validators.ServiceDomainAddUpdate,
	serviceNodes *[]validators.ServiceNodeAddUpdate) error {

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

	// @todo 选择 请求协议： HTTP 和 HTTP&HTTPS 时校验证书是否存在

	createErr := serviceModel.ServiceAdd(&createServiceData, &serviceDomainInfos, &serviceNodeInfos)

	// @todo 记录错误信息的日志，并返回定义的业务提示错误信息

	return createErr
}

func ServiceUpdate(
	serviceId string,
	serviceData *validators.ServiceAddUpdate,
	serviceDomains *[]validators.ServiceDomainAddUpdate,
	serviceNodes *[]validators.ServiceNodeAddUpdate) error {

	updateServiceData := models.Services{
		Protocol:    serviceData.Protocol,
		HealthCheck: serviceData.HealthCheck,
		WebSocket:   serviceData.WebSocket,
		IsEnable:    serviceData.IsEnable,
		LoadBalance: serviceData.LoadBalance,
		Timeouts:    serviceData.Timeouts,
	}

	addDomains, deleteDomainIds := GetToOperateDomains(serviceId, serviceDomains)
	addNodes, updateNodes, deleteNodeIds := GetToOperateNodes(serviceId, serviceNodes)

	// @todo 选择 请求协议： HTTP 和 HTTP&HTTPS 时校验证书是否存在

	serviceModel := &models.Services{}
	updateErr := serviceModel.ServiceUpdate(serviceId, &updateServiceData, &addDomains, &addNodes, &updateNodes, deleteDomainIds, deleteNodeIds)

	// @todo 记录错误信息的日志，并返回定义的业务提示错误信息

	return updateErr
}

func ServiceListPage(param *validators.ServiceList) ([]models.Services, int, error) {
	serviceModel := models.Services{}
	searchContent := strings.TrimSpace(param.Search)

	serviceIds := []string{}
	var listError error
	if len(searchContent) != 0 {
		serviceInfos, serviceErr := serviceModel.ServiceInfosLikeIdName(searchContent)
		if serviceErr != nil {
			listError = serviceErr
		}

		serviceDomainModel := models.ServiceDomains{}
		serviceDomains, domainErr:= serviceDomainModel.ServiceDomainInfosLikeDomain(searchContent)
		if domainErr != nil {
			listError = domainErr
		}

		tpmServiceIds := map[string]string{}
		if len(serviceInfos) != 0 {
			for _, serviceInfo := range serviceInfos {
				_, serviceExist := tpmServiceIds[serviceInfo.ID]
				if !serviceExist {
					tpmServiceIds[serviceInfo.ID] = serviceInfo.ID
				}
			}
		}
		if len(serviceDomains) != 0 {
			for _, serviceDomain := range serviceDomains {
				_, domainExist := tpmServiceIds[serviceDomain.ServiceID]
				if !domainExist {
					tpmServiceIds[serviceDomain.ServiceID] = serviceDomain.ServiceID
				}
			}
		}

		if len(tpmServiceIds) > 0 {
			for _, tpmServiceId := range tpmServiceIds {
				serviceIds = append(serviceIds, tpmServiceId)
			}
		}

		if len(serviceIds) == 0 {
			serviceIds = append(serviceIds, "search-content-exist-set-default-service-id")
		}
	}
	list, total, listError := serviceModel.ServiceAllInfosListPage(serviceIds, param)

	return list, total, listError
}
