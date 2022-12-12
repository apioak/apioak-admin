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
	"time"
)

func CheckServiceExist(serviceResId string) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceResId)
	if len(serviceInfo.ResID) == 0 {
		return errors.New(enums.CodeMessages(enums.ServiceNull))
	}

	return nil
}

func CheckServiceEnableChange(serviceId string, enable int) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)
	if serviceInfo.ResID != serviceId {
		return errors.New(enums.CodeMessages(enums.ServiceNull))
	}

	if serviceInfo.Enable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

//func CheckServiceWebsocketChange(serviceId string, webSocket int) error {
//	serviceModel := &models.Services{}
//	serviceInfo := serviceModel.ServiceInfoById(serviceId)
//	if serviceInfo.ID != serviceId {
//		return errors.New(enums.CodeMessages(enums.ServiceNull))
//	}
//
//	if serviceInfo.WebSocket == webSocket {
//		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
//	}
//
//	return nil
//}
//
//func CheckServiceHealthCheckChange(serviceId string, healthCheck int) error {
//	serviceModel := &models.Services{}
//	serviceInfo := serviceModel.ServiceInfoById(serviceId)
//	if serviceInfo.ID != serviceId {
//		return errors.New(enums.CodeMessages(enums.ServiceNull))
//	}
//
//	if serviceInfo.HealthCheck == healthCheck {
//		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
//	}
//
//	return nil
//}

func CheckServiceRelease(serviceId string) error {
	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)

	if serviceInfo.Release == utils.ReleaseStatusY {
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

	if serviceInfo.Release == utils.ReleaseStatusY {
		if serviceInfo.Enable == utils.EnableOn {
			return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
		}
	} else if serviceInfo.Release == utils.ReleaseStatusT {
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

func ServiceCreate(request *validators.ServiceAddUpdate) error {
	serviceModel := &models.Services{}

	createServiceData := models.Services{
		Name:     request.Name,
		Protocol: request.Protocol,
		Enable:   request.Enable,
		Release:  utils.ReleaseStatusU,
	}

	_, err := serviceModel.ServiceAdd(&createServiceData, request.ServiceDomains)

	if err != nil {
		return err
	}

	return nil
}

func ServiceUpdate(serviceId string, serviceData *validators.ServiceAddUpdate) error {
	serviceModel := models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(serviceId)

	updateServiceData := models.Services{
		Protocol: serviceData.Protocol,
		Enable:   serviceData.Enable,
		Release:  serviceInfo.Release,
	}
	if serviceInfo.Release == utils.ReleaseStatusY {
		updateServiceData.Release = utils.ReleaseStatusT
	}

	if serviceData.Release == utils.ReleaseY {
		updateServiceData.Release = utils.ReleaseStatusY
	}

	serviceDomains := make([]validators.ServiceDomainAddUpdate, 0)
	for _, domain := range serviceData.ServiceDomains {
		serviceDomain := validators.ServiceDomainAddUpdate{
			Domain: domain,
		}

		serviceDomains = append(serviceDomains, serviceDomain)
	}

	addDomains, _ := GetToOperateDomains(serviceId, &serviceDomains)

	updateErr := serviceModel.ServiceUpdate(serviceId, &updateServiceData, &addDomains)

	if (updateErr == nil) && (serviceData.Release == utils.ReleaseY) {
		releaseErr := ServiceRelease(serviceId)
		if releaseErr != nil {
			if serviceInfo.Release != utils.ReleaseStatusU {
				updateServiceData.Release = utils.ReleaseStatusT
			}
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

	// 获取该服务下所有的已发布和待发布的路由和路由插件
	routeModel := models.Routes{}
	releaseRouteInfos := routeModel.RouteInfosByServiceIdReleaseStatus(serviceId, []int{})

	routeIds := make([]string, 0)
	if len(releaseRouteInfos) != 0 {
		for _, releaseRouteInfo := range releaseRouteInfos {
			routeIds = append(routeIds, releaseRouteInfo.ResID)
		}
	}

	routePluginModel := models.RoutePlugins{}
	routePluginInfos := routePluginModel.RoutePluginInfosByRouteIdRelease(routeIds, []int{utils.ReleaseStatusT, utils.ReleaseStatusY})

	if len(releaseRouteInfos) != 0 {
		for _, releaseRouteInfo := range releaseRouteInfos {
			if releaseRouteInfo.Release != utils.ReleaseStatusU {
				routeDeleteReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypeDelete, releaseRouteInfo.ResID)
				if routeDeleteReleaseErr != nil {
					return routeDeleteReleaseErr
				}
			}
		}
	}

	if len(routePluginInfos) != 0 {
		for _, routePluginInfo := range routePluginInfos {
			routePluginDeleteReleaseErr := ServiceRoutePluginConfigRelease(utils.ReleaseTypeDelete, routePluginInfo.ID)
			if routePluginDeleteReleaseErr != nil {
				return routePluginDeleteReleaseErr
			}
		}
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
	ID             string   `json:"id"`              //Service id
	Name           string   `json:"name"`            //Service name
	Protocol       int      `json:"protocol"`        //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	Enable         int      `json:"enable"`          //Service enable  1:on  2:off
	Release        int      `json:"release"`         //Service release status 1:unpublished  2:to be published  3:published
	ServiceDomains []string `json:"service_domains"` //Domain name
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
				_, serviceExist := tpmServiceIds[serviceInfo.ResID]
				if !serviceExist {
					tpmServiceIds[serviceInfo.ResID] = serviceInfo.ResID
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
			tmpServiceInfo.ID = serviceInfo.ResID
			tmpServiceInfo.Name = serviceInfo.Name
			tmpServiceInfo.Protocol = serviceInfo.Protocol
			tmpServiceInfo.Enable = serviceInfo.Enable
			tmpServiceInfo.Release = serviceInfo.Release

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
	ID             string   `json:"id"`              //Service id
	Name           string   `json:"name"`            //Service name
	Protocol       int      `json:"protocol"`        //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	Enable         int      `json:"enable"`          //Service enable  1:on  2:off
	Release        int      `json:"release"`         //Service release status 1:unpublished  2:to be published  3:published
	ServiceDomains []string `json:"service_domains"` //Service Domains
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
	serviceInfo.ID = serviceListInfo.ResID
	serviceInfo.Name = serviceListInfo.Name
	serviceInfo.Protocol = serviceListInfo.Protocol
	serviceInfo.Enable = serviceListInfo.Enable
	serviceInfo.Release = serviceListInfo.Release

	serviceInfo.ServiceDomains = make([]string, 0)
	if len(serviceListInfo.Domains) != 0 {
		for _, domainInfo := range serviceListInfo.Domains {
			serviceInfo.ServiceDomains = append(serviceInfo.ServiceDomains, domainInfo.Domain)
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
	if configReleaseErr == nil {
		routeModel := models.Routes{}
		defaultRouteInfos, defaultRouteInfoErr := routeModel.RouteInfosByServiceRoutePath(serviceId, []string{utils.DefaultRoutePath}, []string{})
		if len(defaultRouteInfos) == 0 {
			serviceModel.ServiceSwitchRelease(serviceId, serviceInfo.Release)
			return errors.New(enums.CodeMessages(enums.RouteDefaultPathNull))
		}

		if defaultRouteInfoErr != nil {
			serviceModel.ServiceSwitchRelease(serviceId, serviceInfo.Release)
			return defaultRouteInfoErr
		}
		defaultRouteInfo := defaultRouteInfos[0]

		routeReleaseErr := ServiceRouteConfigRelease(utils.ReleaseTypePush, defaultRouteInfo.ResID)
		if routeReleaseErr != nil {
			serviceModel.ServiceSwitchRelease(serviceId, serviceInfo.Release)
			return routeReleaseErr
		}

		routeModel.Release = utils.ReleaseStatusY
		routeModel.RouteUpdate(defaultRouteInfo.ResID, routeModel)
	}

	return configReleaseErr
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*utils.EtcdTimeOut)
	defer cancel()

	var respErr error
	if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Put(ctx, etcdKey, serviceConfigStr)
	} else if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Delete(ctx, etcdKey)
	}

	if respErr != nil {
		return errors.New(enums.CodeMessages(enums.EtcdUnavailable))
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
	ID         string            `json:"id"`
	Protocol   int               `json:"protocol"`
	Enable     int               `json:"enable"`
	Domains    []string          `json:"domains"`
	DomainSnis map[string]string `json:"domain_snis"`
}

func generateServicesConfig(id string) (ServiceConfig, error) {
	serviceConfig := ServiceConfig{}
	serviceModel := models.Services{}
	serviceInfo, serviceInfoErr := serviceModel.ServiceDomainNodeById(id)
	if serviceInfoErr != nil {
		return serviceConfig, serviceInfoErr
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

	serviceConfig.ID = serviceInfo.ResID
	serviceConfig.Protocol = serviceInfo.Protocol
	serviceConfig.Enable = serviceInfo.Enable
	serviceConfig.Domains = domains
	serviceConfig.DomainSnis = domainSnis

	return serviceConfig, nil
}
