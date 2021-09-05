package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"encoding/json"
	"errors"
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

type structTimeouts struct {
	ConnectionTimeout int `json:"connection_timeout"`
	ReadTimeout       int `json:"read_timeout"`
	SendTimeout       int `json:"send_timeout"`
}

type StructServiceList struct {
	ID          string         `json:"id"`           //Service id
	Name        string         `json:"name"`         //Service name
	Protocol    int            `json:"protocol"`     //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	HealthCheck int            `json:"health_check"` //Health check switch  1:on  2:off
	WebSocket   int            `json:"web_socket"`   //WebSocket  1:on  2:off
	IsEnable    int            `json:"is_enable"`    //Service enable  1:on  2:off
	LoadBalance int            `json:"load_balance"` //Load balancing algorithm
	Timeouts    structTimeouts `json:"timeouts"`     //Time out
	DomainList  []string       `json:"domain_list"`  //Domain name
}

func (structServiceList *StructServiceList) ServiceListPage(param *validators.ServiceList) ([]StructServiceList, int, error) {
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
		serviceDomains, domainErr := serviceDomainModel.ServiceDomainInfosLikeDomain(searchContent)
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

	serviceList := []StructServiceList{}
	if len(list) != 0 {
		for _, serviceInfo := range list {
			tmpServiceInfo := StructServiceList{}
			tmpServiceInfo.ID = serviceInfo.ID
			tmpServiceInfo.Name = serviceInfo.Name
			tmpServiceInfo.Protocol = serviceInfo.Protocol
			tmpServiceInfo.HealthCheck = serviceInfo.HealthCheck
			tmpServiceInfo.WebSocket = serviceInfo.WebSocket
			tmpServiceInfo.IsEnable = serviceInfo.IsEnable
			tmpServiceInfo.LoadBalance = serviceInfo.LoadBalance

			tmpTimeOuts := structTimeouts{}
			tmpServiceInfo.Timeouts = tmpTimeOuts
			if len(serviceInfo.Timeouts) != 0 {
				tmpTimeOutsErr := json.Unmarshal([]byte(serviceInfo.Timeouts), &tmpTimeOuts)
				if tmpTimeOutsErr == nil {
					tmpServiceInfo.Timeouts = tmpTimeOuts
				}
			}

			tmpServiceInfo.DomainList = make([]string, 0)
			if len(serviceInfo.Domains) != 0 {
				for _, domainInfo := range serviceInfo.Domains {
					tmpServiceInfo.DomainList = append(tmpServiceInfo.DomainList, domainInfo.Domain)
				}
			}

			serviceList = append(serviceList, tmpServiceInfo)
		}
	}

	return serviceList, total, listError
}

type structServiceNode struct {
	NodeIP     string `json:"node_ip"`     //Node IP
	NodePort   int    `json:"node_port"`   //Node port
	NodeWeight int    `json:"node_weight"` //Node weight
}

type StructServiceInfo struct {
	ID          string              `json:"id"`           //Service id
	Name        string              `json:"name"`         //Service name
	Protocol    int                 `json:"protocol"`     //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	HealthCheck int                 `json:"health_check"` //Health check switch  1:on  2:off
	WebSocket   int                 `json:"web_socket"`   //WebSocket  1:on  2:off
	IsEnable    int                 `json:"is_enable"`    //Service enable  1:on  2:off
	LoadBalance int                 `json:"load_balance"` //Load balancing algorithm
	Timeouts    structTimeouts      `json:"timeouts"`     //Time out
	DomainList  []string            `json:"domain_list"`  //Domain name
	NodeList    []structServiceNode `json:"node_list"`    //Domain name
}

func (s *StructServiceInfo) ServiceInfoById(serviceId string) (StructServiceInfo, error) {
	serviceInfo := StructServiceInfo{}
	serviceId = strings.TrimSpace(serviceId)
	err := errors.New(enums.CodeMessages(enums.ServiceParamsNull))
	if len(serviceId) == 0 {
		return serviceInfo, err
	}

	serviceModel := models.Services{}
	serviceList, err := serviceModel.ServiceDomainNodeByIds([]string{serviceId})
	if err != nil {
		return serviceInfo, err
	}

	serviceListInfo := serviceList[0]
	serviceInfo.ID = serviceListInfo.ID
	serviceInfo.Name = serviceListInfo.Name
	serviceInfo.Protocol = serviceListInfo.Protocol
	serviceInfo.HealthCheck = serviceListInfo.HealthCheck
	serviceInfo.WebSocket = serviceListInfo.WebSocket
	serviceInfo.IsEnable = serviceListInfo.IsEnable
	serviceInfo.LoadBalance = serviceListInfo.LoadBalance

	tmpTimeOuts := structTimeouts{}
	serviceInfo.Timeouts = tmpTimeOuts
	if len(serviceListInfo.Timeouts) != 0 {
		tmpTimeOutsErr := json.Unmarshal([]byte(serviceListInfo.Timeouts), &tmpTimeOuts)
		if tmpTimeOutsErr == nil {
			serviceInfo.Timeouts = tmpTimeOuts
		}
	}

	serviceInfo.DomainList = make([]string, 0)
	if len(serviceListInfo.Domains) != 0 {
		for _, domainInfo := range serviceListInfo.Domains {
			serviceInfo.DomainList = append(serviceInfo.DomainList, domainInfo.Domain)
		}
	}

	serviceInfo.NodeList = make([]structServiceNode, 0)
	if len(serviceListInfo.Nodes) != 0 {
		for _, nodeInfo := range serviceListInfo.Nodes {
			tmpNodeInfo := structServiceNode{}
			tmpNodeInfo.NodeIP = nodeInfo.NodeIP
			tmpNodeInfo.NodePort = nodeInfo.NodePort
			tmpNodeInfo.NodeWeight = nodeInfo.NodeWeight

			serviceInfo.NodeList = append(serviceInfo.NodeList, tmpNodeInfo)
		}
	}

	return serviceInfo, nil
}
