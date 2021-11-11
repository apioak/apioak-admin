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

func CheckServiceExist(serviceId string) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		return errors.New(enums.CodeMessages(enums.ServiceNull))
	}

	return nil
}

func CheckExistDomain(domains []string, filterServiceIds []string) error {
	serviceDomainInfo := models.ServiceDomains{}
	serviceDomains, err := serviceDomainInfo.DomainInfosByDomain(domains, filterServiceIds)
	if err != nil {
		return nil
	}

	if len(serviceDomains) == 0 {
		return nil
	}

	existDomains := make([]string, 0)
	tmpExistDomainsMap := make(map[string]byte, 0)
	for _, serviceDomain := range serviceDomains {
		_, exist := tmpExistDomainsMap[serviceDomain.Domain]
		if exist {
			continue
		}

		existDomains = append(existDomains, serviceDomain.Domain)
		tmpExistDomainsMap[serviceDomain.Domain] = 0
	}

	if len(existDomains) != 0 {
		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.ServiceDomainExist), strings.Join(existDomains, ",")))
	}

	return nil
}

func CheckDomainCertificate(protocol int, domains []string) error {
	if (protocol != utils.ProtocolHTTPS) && (protocol != utils.ProtocolHTTPAndHTTPS) {
		return nil
	}

	domainSniInfos, domainSniInfosErr := utils.InterceptSni(domains)
	if domainSniInfosErr != nil {
		return domainSniInfosErr
	}

	certificatesModel := models.Certificates{}
	domainCertificateInfos := certificatesModel.CertificateInfoByDomainSniInfos(domainSniInfos)
	if len(domainCertificateInfos) == len(domainSniInfos) {
		return nil
	}

	nullCertificateDomains := make([]string, 0)
	for _, domainInfo := range domains {
		if len(domainCertificateInfos) == 0 {

			nullCertificateDomains = append(nullCertificateDomains, domainInfo)
		} else {
			for _, domainCertificateInfo := range domainCertificateInfos {

				domainSni := strings.ReplaceAll(domainCertificateInfo.Sni, "*", "")
				if domainInfo[len(domainInfo)-len(domainSni):] != domainSni {
					nullCertificateDomains = append(nullCertificateDomains, domainInfo)
				}
			}
		}
	}

	if len(nullCertificateDomains) != 0 {
		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.ServiceDomainSslNull), strings.Join(nullCertificateDomains, ",")))
	}

	return nil
}

func ServiceCreate(serviceData *validators.ServiceAddUpdate) error {
	serviceModel := &models.Services{}
	serviceDomainInfos := make([]models.ServiceDomains, 0)
	serviceNodeInfos := make([]models.ServiceNodes, 0)

	timeOutByte, _ := json.Marshal(serviceData.Timeouts)
	createServiceData := models.Services{
		Protocol:    serviceData.Protocol,
		HealthCheck: serviceData.HealthCheck,
		WebSocket:   serviceData.WebSocket,
		IsEnable:    serviceData.IsEnable,
		LoadBalance: serviceData.LoadBalance,
		Timeouts:    string(timeOutByte),
	}

	for _, domainInfo := range serviceData.ServiceDomains {
		domain := models.ServiceDomains{
			Domain: domainInfo,
		}
		serviceDomainInfos = append(serviceDomainInfos, domain)
	}

	for _, nodeInfo := range serviceData.ServiceNodes {
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

	// @todo 如果状态是"开启"，则需要同步远程数据中心

	// @todo 记录错误信息的日志，并返回定义的业务提示错误信息

	return createErr
}

func ServiceUpdate(serviceId string, serviceData *validators.ServiceAddUpdate) error {

	timeOutByte, _ := json.Marshal(serviceData.Timeouts)
	updateServiceData := models.Services{
		Protocol:    serviceData.Protocol,
		HealthCheck: serviceData.HealthCheck,
		WebSocket:   serviceData.WebSocket,
		IsEnable:    serviceData.IsEnable,
		LoadBalance: serviceData.LoadBalance,
		Timeouts:    string(timeOutByte),
	}

	serviceDomains := make([]validators.ServiceDomainAddUpdate, 0)
	for _, domain := range serviceData.ServiceDomains {
		serviceDomain := validators.ServiceDomainAddUpdate{
			Domain: domain,
		}

		serviceDomains = append(serviceDomains, serviceDomain)
	}

	addDomains, deleteDomainIds := GetToOperateDomains(serviceId, &serviceDomains)
	addNodes, updateNodes, deleteNodeIds := GetToOperateNodes(serviceId, &serviceData.ServiceNodes)

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
	ID             string         `json:"id"`              //Service id
	Name           string         `json:"name"`            //Service name
	Protocol       int            `json:"protocol"`        //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	HealthCheck    int            `json:"health_check"`    //Health check switch  1:on  2:off
	WebSocket      int            `json:"web_socket"`      //WebSocket  1:on  2:off
	IsEnable       int            `json:"is_enable"`       //Service enable  1:on  2:off
	IsRelease      int            `json:"is_release"`      //Service release  1:on  2:off
	LoadBalance    int            `json:"load_balance"`    //Load balancing algorithm
	Timeouts       structTimeouts `json:"timeouts"`        //Time out
	ServiceDomains []string       `json:"service_domains"` //Domain name
}

func (structServiceList *StructServiceList) ServiceListPage(param *validators.ServiceList) ([]StructServiceList, int, error) {
	serviceModel := models.Services{}
	searchContent := strings.TrimSpace(param.Search)

	serviceIds := make([]string, 0)
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

	serviceList := make([]StructServiceList, 0)
	if len(list) != 0 {
		for _, serviceInfo := range list {
			tmpServiceInfo := StructServiceList{}
			tmpServiceInfo.ID = serviceInfo.ID
			tmpServiceInfo.Name = serviceInfo.Name
			tmpServiceInfo.Protocol = serviceInfo.Protocol
			tmpServiceInfo.HealthCheck = serviceInfo.HealthCheck
			tmpServiceInfo.WebSocket = serviceInfo.WebSocket
			tmpServiceInfo.IsEnable = serviceInfo.IsEnable
			tmpServiceInfo.IsRelease = serviceInfo.IsRelease
			tmpServiceInfo.LoadBalance = serviceInfo.LoadBalance

			tmpTimeOuts := structTimeouts{}
			tmpServiceInfo.Timeouts = tmpTimeOuts
			if len(serviceInfo.Timeouts) != 0 {
				tmpTimeOutsErr := json.Unmarshal([]byte(serviceInfo.Timeouts), &tmpTimeOuts)
				if tmpTimeOutsErr == nil {
					tmpServiceInfo.Timeouts = tmpTimeOuts
				}
			}

			tmpServiceInfo.ServiceDomains = make([]string, 0)
			if len(serviceInfo.Domains) != 0 {
				for _, domainInfo := range serviceInfo.Domains {
					tmpServiceInfo.ServiceDomains = append(tmpServiceInfo.ServiceDomains, domainInfo.Domain)
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
	ID             string              `json:"id"`              //Service id
	Name           string              `json:"name"`            //Service name
	Protocol       int                 `json:"protocol"`        //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	HealthCheck    int                 `json:"health_check"`    //Health check switch  1:on  2:off
	WebSocket      int                 `json:"web_socket"`      //WebSocket  1:on  2:off
	IsEnable       int                 `json:"is_enable"`       //Service enable  1:on  2:off
	LoadBalance    int                 `json:"load_balance"`    //Load balancing algorithm
	Timeouts       structTimeouts      `json:"timeouts"`        //Time out
	ServiceDomains []string            `json:"service_domains"` //Service Domains
	ServiceNodes   []structServiceNode `json:"service_nodes"`   //Service Nodes
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

	serviceInfo.ServiceDomains = make([]string, 0)
	if len(serviceListInfo.Domains) != 0 {
		for _, domainInfo := range serviceListInfo.Domains {
			serviceInfo.ServiceDomains = append(serviceInfo.ServiceDomains, domainInfo.Domain)
		}
	}

	serviceInfo.ServiceNodes = make([]structServiceNode, 0)
	if len(serviceListInfo.Nodes) != 0 {
		for _, nodeInfo := range serviceListInfo.Nodes {
			tmpNodeInfo := structServiceNode{}
			tmpNodeInfo.NodeIP = nodeInfo.NodeIP
			tmpNodeInfo.NodePort = nodeInfo.NodePort
			tmpNodeInfo.NodeWeight = nodeInfo.NodeWeight

			serviceInfo.ServiceNodes = append(serviceInfo.ServiceNodes, tmpNodeInfo)
		}
	}

	return serviceInfo, nil
}
