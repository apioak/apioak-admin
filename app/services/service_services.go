package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type ServicesService struct {
}

var (
	servicesService *ServicesService
	servicesOnce    sync.Once
)

func NewServicesService() *ServicesService {

	servicesOnce.Do(func() {
		servicesService = &ServicesService{}
	})

	return servicesService
}

func CheckServiceExist(serviceResId string) error {
	serviceModel := &models.Services{}
	_, err := serviceModel.ServiceInfoById(serviceResId)

	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) ServiceCreate(request *validators.ServiceAddUpdate) error {

	createServiceData := &models.Services{
		Name:     request.Name,
		Protocol: request.Protocol,
		Enable:   request.Enable,
		Release:  utils.ReleaseStatusU,
	}

	_, err := (&models.Services{}).ServiceAdd(createServiceData, request.ServiceDomains)

	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) ServiceUpdate(serviceId string, request *validators.ServiceAddUpdate) error {
	serviceModel := models.Services{}
	serviceInfo, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	updateServiceData := models.Services{
		Name:     request.Name,
		Protocol: request.Protocol,
		Enable:   request.Enable,
		Release:  serviceInfo.Release,
	}
	if serviceInfo.Release == utils.ReleaseStatusY {
		updateServiceData.Release = utils.ReleaseStatusT
	}

	return serviceModel.ServiceUpdate(serviceId, &updateServiceData, request.ServiceDomains)
}

func (s *ServicesService) ServiceUpdateName(serviceId string, request *validators.ServiceUpdateName) error {
	serviceModel := models.Services{}
	service, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	updateParam := map[string]interface{}{
		"name": request.Name,
	}
	if service.Release == utils.ReleaseStatusY {
		updateParam["release"] = utils.ReleaseStatusT
	}

	return serviceModel.ServiceUpdateColumns(serviceId, updateParam)
}

func checkServiceEnableChange(serviceId string, enable int) error {
	serviceModel := &models.Services{}
	serviceInfo, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	if serviceInfo.Enable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func (s *ServicesService) ServiceSwitchEnable(serviceId string, enable int) error {
	serviceModel := models.Services{}
	service, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	err = checkServiceEnableChange(serviceId, enable)

	updateParam := map[string]interface{}{
		"enable": enable,
	}
	if service.Release == utils.ReleaseStatusY {
		updateParam["release"] = utils.ReleaseStatusT
	}

	return serviceModel.ServiceUpdateColumns(serviceId, updateParam)
}

type StructServiceInfo struct {
	ID             int64    `json:"id"`              //Service id
	ResID          string   `json:"res_id"`          //Service res id
	Name           string   `json:"name"`            //Service name
	Protocol       int      `json:"protocol"`        //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	Enable         int      `json:"enable"`          //Service enable  1:on  2:off
	Release        int      `json:"release"`         //Service release status 1:unpublished  2:to be published  3:published
	ServiceDomains []string `json:"service_domains"` //Service Domains
}

func (s *ServicesService) ServiceInfoById(serviceId string) (StructServiceInfo, error) {
	serviceInfo := StructServiceInfo{}

	service, err := (&models.Services{}).ServiceInfoById(serviceId)
	if err != nil {
		return serviceInfo, err
	}

	serviceInfo = StructServiceInfo{
		ID:       service.ID,
		ResID:    service.ResID,
		Name:     service.Name,
		Protocol: service.Protocol,
		Enable:   service.Enable,
		Release:  service.Release,
	}
	serviceDomain, err := (&models.ServiceDomains{}).DomainInfosByServiceIds([]string{serviceId})

	domain := []string{}
	if err == nil {
		for _, v := range serviceDomain {
			domain = append(domain, v.Domain)
		}
	}

	serviceInfo.ServiceDomains = domain

	return serviceInfo, nil
}

func (s *ServicesService) ServiceDelete(serviceId string) error {

	routeModel := models.Routes{}
	routerList, err := routeModel.RouteInfosByServiceId(serviceId)

	if err != nil {
		return errors.New(err.Error())
	}

	if len(routerList) > 0 {
		return errors.New(enums.CodeMessages(enums.ServiceBindingRouter))
	}

	// TODO 获取consul 服务数据

	serviceModel := &models.Services{}
	err = serviceModel.ServiceDelete(serviceId)
	if err != nil {
		return errors.New(err.Error())
	}

	// TODO 删除consul 服务数据

	return nil
}

type ServiceItem struct {
	ID             int64    `json:"id"`
	ResID          string   `json:"res_id"`          //Service id
	Name           string   `json:"name"`            //Service name
	Protocol       int      `json:"protocol"`        //Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS
	Enable         int      `json:"enable"`          //Service enable  1:on  2:off
	Release        int      `json:"release"`         //Service release status 1:unpublished  2:to be published  3:published
	ServiceDomains []string `json:"service_domains"` //Domain name
}

func (s *ServicesService) ServiceList(request *validators.ServiceList) ([]ServiceItem, int, error) {
	serviceModel := models.Services{}
	searchContent := strings.TrimSpace(request.Search)

	serviceIds := []string{}
	if searchContent != "" {
		services, err := serviceModel.ServiceInfosLikeResIdName(searchContent)
		if err == nil {
			for _, v := range services {
				serviceIds = append(serviceIds, v.ResID)
			}
		}

		serviceDomainModel := models.ServiceDomains{}
		serviceDomains, err := serviceDomainModel.ServiceDomainInfosLikeDomain(searchContent)
		if err == nil {
			for _, v := range serviceDomains {
				serviceIds = append(serviceIds, v.ServiceResID)
			}
		}

		if len(serviceIds) == 0 {
			return []ServiceItem{}, 0, nil
		}
	}
	list, total, err := serviceModel.ServiceList(serviceIds, request)

	if err != nil && err != gorm.ErrRecordNotFound {
		return []ServiceItem{}, 0, err
	}

	serviceList := []ServiceItem{}

	if len(list) == 0 {
		return []ServiceItem{}, 0, nil
	}

	listServiceId := []string{}
	for _, v := range list {
		listServiceId = append(listServiceId, v.ResID)
	}

	domains, err := (&models.ServiceDomains{}).DomainInfosByServiceIds(listServiceId)

	serviceDomainMap := map[string][]models.ServiceDomains{}
	if err == nil {
		for _, v := range domains {
			if _, ok := serviceDomainMap[v.ServiceResID]; !ok {
				serviceDomainMap[v.ServiceResID] = []models.ServiceDomains{}
			}
			serviceDomainMap[v.ServiceResID] = append(serviceDomainMap[v.ServiceResID], v)
		}
	}

	for _, v := range list {

		domain := []string{}
		if tmp, ok := serviceDomainMap[v.ResID]; ok {
			for _, vd := range tmp {
				domain = append(domain, vd.Domain)
			}
		}
		serviceList = append(serviceList, ServiceItem{
			ID:             v.ID,
			ResID:          v.ResID,
			Name:           v.Name,
			Protocol:       v.Protocol,
			Enable:         v.Enable,
			Release:        v.Release,
			ServiceDomains: domain,
		})
	}

	return serviceList, total, nil
}

func checkServiceRelease(serviceId string) error {
	serviceModel := &models.Services{}
	serviceInfo, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}

	if serviceInfo.Release == utils.ReleaseStatusY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	if (serviceInfo.Protocol == utils.ProtocolHTTPS) || (serviceInfo.Protocol == utils.ProtocolHTTPAndHTTPS) {
		serviceDomains, err := (&models.ServiceDomains{}).DomainInfosByServiceIds([]string{serviceId})
		if err != nil {
			return err
		}
		if len(serviceDomains) == 0 {
			return nil
		}

		domains := []string{}
		for _, v := range serviceDomains {
			domains = append(domains, v.Domain)
		}

		domainSnis, err := utils.InterceptSni(domains)
		if err != nil {
			return err
		}

		certificatesModel := models.Certificates{}
		_ = certificatesModel.CertificateInfoByDomainSniInfos(domainSnis)

		//if len(domainCert) < len(domainSnis) {
		//
		//	domainCertificatesMap := make(map[string]byte, 0)
		//	for _, domainCertificateInfo := range domainCert {
		//		domainCertificatesMap[domainCertificateInfo.Sni] = 0
		//	}
		//
		//	noCertificateDomains := make([]string, 0)
		//	for _, serviceDomainInfo := range domainCert {
		//		disassembleDomains := strings.Split(serviceDomainInfo.Domain, ".")
		//		disassembleDomains[0] = "*"
		//		domainSni := strings.Join(disassembleDomains, ".")
		//		_, ok := domainCertificatesMap[domainSni]
		//		if ok == false {
		//			noCertificateDomains = append(noCertificateDomains, serviceDomainInfo.Domain)
		//		}
		//	}
		//
		//	if len(noCertificateDomains) != 0 {
		//		return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.ServiceDomainSslNull), strings.Join(noCertificateDomains, ",")))
		//	}
		//}
		//
		//noReleaseCertificates := make([]string, 0)
		//noEnableCertificates := make([]string, 0)
		//for _, domainCertificateInfo := range domainCertificateInfos {
		//	if domainCertificateInfo.ReleaseStatus != utils.ReleaseStatusY {
		//		noReleaseCertificates = append(noReleaseCertificates, domainCertificateInfo.Sni)
		//	}
		//	if domainCertificateInfo.IsEnable != utils.EnableOn {
		//		noEnableCertificates = append(noEnableCertificates, domainCertificateInfo.Sni)
		//	}
		//}
		//
		//if len(noReleaseCertificates) != 0 {
		//	return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.CertificateNoRelease), strings.Join(noReleaseCertificates, ",")))
		//}
		//
		//if len(noEnableCertificates) != 0 {
		//	return fmt.Errorf(fmt.Sprintf(enums.CodeMessages(enums.CertificateEnableOff), strings.Join(noEnableCertificates, ",")))
		//}
	}

	return nil
}

func (s *ServicesService) ServiceRelease(serviceId string) error {
	serviceModel := &models.Services{}
	_, err := serviceModel.ServiceInfoById(serviceId)

	if err != nil {
		return err
	}
	err = checkServiceRelease(serviceId)
	if err != nil {
		return err
	}

	updateParam := map[string]interface{}{
		"release": utils.ReleaseStatusY,
	}
	err = serviceModel.ServiceUpdateColumns(serviceId, updateParam)
	if err != nil {
		return err
	}

	// TODO 查询consul service info

	// TODO 更新consul service info
	return nil
}

func (s *ServicesService) CheckExistDomain(domains []string, filterServiceIds []string) error {
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

	domainSniInfos, err := utils.InterceptSni(domains)
	if err != nil {
		return err
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
