package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"context"
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

func CheckServiceEnableChange(serviceId string, enable int) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		return errors.New(enums.CodeMessages(enums.ServiceNull))
	}

	if serviceInfo.IsEnable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func CheckServiceWebsocketChange(serviceId string, webSocket int) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		return errors.New(enums.CodeMessages(enums.ServiceNull))
	}

	if serviceInfo.WebSocket == webSocket {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func CheckServiceHealthCheckChange(serviceId string, healthCheck int) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ID != serviceId {
		return errors.New(enums.CodeMessages(enums.ServiceNull))
	}

	if serviceInfo.HealthCheck == healthCheck {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func CheckServiceRelease(serviceId string) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)

	if serviceInfo.ReleaseStatus == utils.ReleaseStatusY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	if (serviceInfo.Protocol == utils.ProtocolHTTPS) || (serviceInfo.Protocol == utils.ProtocolHTTPAndHTTPS) {
		serviceDomainModel := models.ServiceDomains{}
		serviceDomainInfos, serviceDomainInfosErr := serviceDomainModel.DomainInfosByServiceIds([]string{serviceId})
		if serviceDomainInfosErr != nil {
			return serviceDomainInfosErr
		}
		if len(serviceDomainInfos) == 0 {
			return nil
		}

		serviceDomains := make([]string, 0)
		for _, serviceDomainInfo := range serviceDomainInfos {
			serviceDomains = append(serviceDomains, serviceDomainInfo.Domain)
		}

		domainSnis, domainSnisErr := utils.InterceptSni(serviceDomains)
		if domainSnisErr != nil {
			return domainSnisErr
		}

		certificatesModel := models.Certificates{}
		domainCertificateInfos := certificatesModel.CertificateInfoByDomainSniInfos(domainSnis)

		if len(domainCertificateInfos) < len(domainSnis) {

			domainCertificatesMap := make(map[string]byte, 0)
			for _, domainCertificateInfo := range domainCertificateInfos {
				domainCertificatesMap[domainCertificateInfo.Sni] = 0
			}

			noCertificateDomains := make([]string, 0)
			for _, serviceDomainInfo := range serviceDomainInfos {
				disassembleDomains := strings.Split(serviceDomainInfo.Domain, ".")
				disassembleDomains[0] = "*"
				domainSni := strings.Join(disassembleDomains, ".")
				_, ok := domainCertificatesMap[domainSni]
				if ok == false {
					noCertificateDomains = append(noCertificateDomains, serviceDomainInfo.Domain)
				}
			}

			if len(noCertificateDomains) != 0 {
				return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.ServiceDomainSslNull), strings.Join(noCertificateDomains, ",")))
			}
		}

		noReleaseCertificates := make([]string, 0)
		noEnableCertificates := make([]string, 0)
		for _, domainCertificateInfo := range domainCertificateInfos {
			if domainCertificateInfo.ReleaseStatus != utils.ReleaseStatusY {
				noReleaseCertificates = append(noReleaseCertificates, domainCertificateInfo.Sni)
			}
			if domainCertificateInfo.IsEnable != utils.EnableOn {
				noEnableCertificates = append(noEnableCertificates, domainCertificateInfo.Sni)
			}
		}

		if len(noReleaseCertificates) != 0 {
			return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.CertificateNoRelease), strings.Join(noReleaseCertificates, ",")))
		}

		if len(noEnableCertificates) != 0 {
			return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.CertificateEnableOff), strings.Join(noEnableCertificates, ",")))
		}
	}

	return nil
}

func CheckServiceDelete(serviceId string) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)

	if serviceInfo.ReleaseStatus == utils.ReleaseStatusY {
		if serviceInfo.IsEnable == utils.EnableOn {
			return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
		}
	} else if serviceInfo.ReleaseStatus == utils.ReleaseStatusT {
		return errors.New(enums.CodeMessages(enums.ToReleaseProhibitsOp))
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
		Protocol:      serviceData.Protocol,
		HealthCheck:   serviceData.HealthCheck,
		WebSocket:     serviceData.WebSocket,
		IsEnable:      serviceData.IsEnable,
		ReleaseStatus: utils.ReleaseStatusU,
		LoadBalance:   serviceData.LoadBalance,
		Timeouts:      string(timeOutByte),
	}

	if serviceData.IsRelease == utils.IsReleaseY {
		createServiceData.ReleaseStatus = utils.ReleaseStatusY
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

	serviceId, createErr := serviceModel.ServiceAdd(&createServiceData, &serviceDomainInfos, &serviceNodeInfos)

	if (createErr == nil) && (serviceData.IsRelease == utils.IsReleaseY) {
		releaseErr := ServiceRelease(serviceId)
		if releaseErr != nil {
			createServiceData.ReleaseStatus = utils.ReleaseStatusU
			serviceModel.ServiceUpdateColumnsById(serviceId, &createServiceData)

			return releaseErr
		}
	}

	return createErr
}

func ServiceUpdate(serviceId string, serviceData *validators.ServiceAddUpdate) error {
	serviceModel := models.Services{}
	timeOutByte, _ := json.Marshal(serviceData.Timeouts)
	serviceInfo := serviceModel.ServiceInfoById(serviceId)

	updateServiceData := models.Services{
		Protocol:      serviceData.Protocol,
		HealthCheck:   serviceData.HealthCheck,
		WebSocket:     serviceData.WebSocket,
		IsEnable:      serviceData.IsEnable,
		ReleaseStatus: serviceInfo.ReleaseStatus,
		LoadBalance:   serviceData.LoadBalance,
		Timeouts:      string(timeOutByte),
	}

	if serviceData.IsRelease == utils.IsReleaseY {
		updateServiceData.ReleaseStatus = utils.ReleaseStatusY
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

	updateErr := serviceModel.ServiceUpdate(serviceId, &updateServiceData, &addDomains, &addNodes, &updateNodes, deleteDomainIds, deleteNodeIds)

	if (updateErr == nil) && (serviceData.IsRelease == utils.IsReleaseY) {
		releaseErr := ServiceRelease(serviceId)
		if releaseErr != nil {
			updateServiceData.ReleaseStatus = utils.ReleaseStatusT
			serviceModel.ServiceUpdateColumnsById(serviceId, &updateServiceData)

			return releaseErr
		}
	}

	return updateErr
}

func ServiceDelete(serviceId string) error {
	configReleaseErr := ServiceConfigRelease(utils.ReleaseTypeDelete, serviceId)
	if configReleaseErr != nil {
		return configReleaseErr
	}

	serviceModel := &models.Services{}
	deleteErr := serviceModel.ServiceDelete(serviceId)
	if deleteErr != nil {
		ServiceConfigRelease(utils.ReleaseTypePush, serviceId)
		return errors.New(deleteErr.Error())
	}

	return nil
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
	ReleaseStatus  int            `json:"release_status"`  //Service release status 1:unpublished  2:to be published  3:published
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
			tmpServiceInfo.ReleaseStatus = serviceInfo.ReleaseStatus
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
	ReleaseStatus  int                 `json:"release_status"`  //Service release status 1:unpublished  2:to be published  3:published
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
	serviceInfo.ReleaseStatus = serviceListInfo.ReleaseStatus
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

func ServiceRelease(serviceId string) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)

	serviceReleaseErr := serviceModel.ServiceSwitchRelease(serviceId, utils.ReleaseStatusY)
	if serviceReleaseErr != nil {
		return serviceReleaseErr
	}

	configReleaseErr := ServiceConfigRelease(utils.ReleaseTypePush, serviceId)
	if configReleaseErr != nil {
		serviceModel.ServiceSwitchRelease(serviceId, serviceInfo.ReleaseStatus)
		return configReleaseErr
	}

	return nil
}

func ServiceConfigRelease(releaseType string, serviceId string) error {
	serviceConfig, serviceConfigErr := generateServicesConfig(serviceId)
	if serviceConfigErr != nil {
		return serviceConfigErr
	}

	serviceConfigJson, serviceConfigJsonErr := json.Marshal(serviceConfig)
	if serviceConfigJsonErr != nil {
		return serviceConfigJsonErr
	}
	serviceConfigStr := string(serviceConfigJson)

	etcdKey := utils.EtcdKey(utils.EtcdKeyTypeService, serviceId)
	if len(etcdKey) == 0 {
		return errors.New(enums.CodeMessages(enums.EtcdKeyNull))
	}

	etcdClient := packages.GetEtcdClient()

	var respErr error
	if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Put(context.Background(), etcdKey, serviceConfigStr)
	} else if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Delete(context.Background(), etcdKey)
	}

	if respErr != nil {
		return respErr
	}

	return nil
}

type ServiceUpstreamConfig struct {
	IPType int    `json:"ip_type"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

type ServiceTimeOutConfig struct {
	ConnectionTimeout int `json:"connection_timeout"`
	ReadTimeout       int `json:"read_timeout"`
	SendTimeout       int `json:"send_timeout"`
}

type ServiceConfig struct {
	ID          string                  `json:"id"`
	Protocol    int                     `json:"protocol"`
	HealthCheck int                     `json:"health_check"`
	WebSocket   int                     `json:"web_socket"`
	IsEnable    int                     `json:"is_enable"`
	LoadBalance int                     `json:"load_balance"`
	TimeOut     ServiceTimeOutConfig    `json:"time_out"`
	Domains     []string                `json:"domains"`
	DomainSnis  map[string]string       `json:"domain_snis"`
	Upstreams   []ServiceUpstreamConfig `json:"upstreams"`
}

func generateServicesConfig(id string) (ServiceConfig, error) {
	serviceConfig := ServiceConfig{}
	serviceModel := models.Services{}
	serviceInfo, serviceInfoErr := serviceModel.ServiceDomainNodeById(id)
	if serviceInfoErr != nil {
		return serviceConfig, serviceInfoErr
	}

	serviceTimeOutConfig := ServiceTimeOutConfig{}
	timeOutInfoErr := json.Unmarshal([]byte(serviceInfo.Timeouts), &serviceTimeOutConfig)
	if timeOutInfoErr != nil {
		return serviceConfig, timeOutInfoErr
	}

	domains := make([]string, 0)
	domainSnis := make(map[string]string, 0)
	if len(serviceInfo.Domains) != 0 {
		for _, domain := range serviceInfo.Domains {
			domains = append(domains, domain.Domain)
			disassembleDomains := strings.Split(domain.Domain, ".")
			disassembleDomains[0] = "*"
			domainSniInfo := strings.Join(disassembleDomains, ".")
			domainSnis[domain.Domain] = domainSniInfo
		}
	}

	upstreams := make([]ServiceUpstreamConfig, 0)
	if len(serviceInfo.Nodes) != 0 {
		for _, nodeInfo := range serviceInfo.Nodes {
			serviceUpstreamConfig := ServiceUpstreamConfig{
				IPType: nodeInfo.IPType,
				IP:     nodeInfo.NodeIP,
				Port:   nodeInfo.NodePort,
				Weight: nodeInfo.NodeWeight,
			}

			upstreams = append(upstreams, serviceUpstreamConfig)
		}
	}

	serviceConfig.ID = serviceInfo.ID
	serviceConfig.Protocol = serviceInfo.Protocol
	serviceConfig.HealthCheck = serviceInfo.HealthCheck
	serviceConfig.WebSocket = serviceInfo.WebSocket
	serviceConfig.IsEnable = serviceInfo.IsEnable
	serviceConfig.LoadBalance = serviceInfo.LoadBalance
	serviceConfig.TimeOut = serviceTimeOutConfig
	serviceConfig.Domains = domains
	serviceConfig.DomainSnis = domainSnis
	serviceConfig.Upstreams = upstreams

	return serviceConfig, nil
}
