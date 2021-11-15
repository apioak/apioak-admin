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

func CheckCertificateNull(id string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	if certificateInfo.ID != id {
		return errors.New(enums.CodeMessages(enums.CertificateNull))
	}

	return nil
}

func CheckCertificateExistBySni(sni string, filterId string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoBySni(sni, filterId)
	if certificateInfo.Sni == sni {
		return errors.New(enums.CodeMessages(enums.CertificateExist))
	}

	return nil
}

func CheckCertificateDomainExist(sni string) error {
	searchSni := strings.Replace(sni, "*.", "", 1)

	serviceDomainsModel := models.ServiceDomains{}
	domainList := serviceDomainsModel.DomainListByLikeSni(searchSni)
	if len(domainList) != 0 {
		return errors.New(enums.CodeMessages(enums.CertificateDomainExist))
	}

	return nil
}

func CheckCertificateDomainExistById(id string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	if certificateInfo.ID != id {
		return errors.New(enums.CodeMessages(enums.CertificateNull))
	}

	checkCertificateDomainExistErr := CheckCertificateDomainExist(certificateInfo.Sni)
	if checkCertificateDomainExistErr != nil {
		return checkCertificateDomainExistErr
	}

	return nil
}

func CheckCertificateDelete(id string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)

	if certificateInfo.IsEnable == utils.EnableOn {
		return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
	}

	if certificateInfo.IsRelease != utils.IsReleaseY {
		return errors.New(enums.CodeMessages(enums.EnablePublishedONOp))
	}

	return nil
}

func CheckCertificateEnableChange(id string, enable int) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	if certificateInfo.IsEnable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func CheckCertificateRelease(id string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	if certificateInfo.IsRelease == utils.IsReleaseY {
		return errors.New(enums.CodeMessages(enums.SwitchPublished))
	}

	return nil
}

func decodeCertificateData(certificateContent string) (string, error) {
	certificateInfo := ""
	type contentStruct struct {
		Content string `json:"content"`
	}

	contentInfo := contentStruct{}
	contentInfoErr := json.Unmarshal([]byte(certificateContent), &contentInfo)
	if contentInfoErr != nil {
		return certificateInfo, contentInfoErr
	}

	certificateInfo = contentInfo.Content

	return certificateInfo, nil
}

func CertificateAdd(certificateData *validators.CertificateAddUpdate) error {
	certificateContent, certificateContentErr := decodeCertificateData(certificateData.Certificate)
	if certificateContentErr != nil {
		return certificateContentErr
	}
	certificateData.Certificate = certificateContent

	privateKeyContent, privateKeyContentErr := decodeCertificateData(certificateData.PrivateKey)
	if privateKeyContentErr != nil {
		return privateKeyContentErr
	}
	certificateData.PrivateKey = privateKeyContent

	certificateInfo, certificateInfoErr := utils.DiscernCertificate(&certificateData.Certificate)
	if certificateInfoErr != nil {
		return certificateInfoErr
	}

	checkCertificateExistErr := CheckCertificateExistBySni(certificateInfo.CommonName, "")
	if checkCertificateExistErr != nil {
		return checkCertificateExistErr
	}

	certificatesModel := models.Certificates{
		Certificate: certificateContent,
		PrivateKey:  privateKeyContent,
		ExpiredAt:   certificateInfo.NotAfter,
		IsEnable:    certificateData.IsEnable,
		IsRelease:   utils.IsReleaseN,
		Sni:         certificateInfo.CommonName,
	}

	addErr := certificatesModel.CertificatesAdd(&certificatesModel)
	if addErr != nil {
		return addErr
	}

	return nil
}

func CertificateUpdate(id string, certificateData *validators.CertificateAddUpdate) error {
	certificateContent, certificateContentErr := decodeCertificateData(certificateData.Certificate)
	if certificateContentErr != nil {
		return certificateContentErr
	}
	certificateData.Certificate = certificateContent

	privateKeyContent, privateKeyContentErr := decodeCertificateData(certificateData.PrivateKey)
	if privateKeyContentErr != nil {
		return privateKeyContentErr
	}
	certificateData.PrivateKey = privateKeyContent

	certificateInfo, certificateInfoErr := utils.DiscernCertificate(&certificateData.Certificate)
	if certificateInfoErr != nil {
		return certificateInfoErr
	}

	certificatesModel := models.Certificates{}
	certificateExistInfo := certificatesModel.CertificateInfoById(id)

	if certificateInfo.CommonName != certificateExistInfo.Sni {
		checkCertificateDomainExistErr := CheckCertificateDomainExist(certificateExistInfo.Sni)
		if checkCertificateDomainExistErr != nil {
			return checkCertificateDomainExistErr
		}
	}

	checkCertificateExistErr := CheckCertificateExistBySni(certificateInfo.CommonName, id)
	if checkCertificateExistErr != nil {
		return checkCertificateExistErr
	}

	certificatesModel.Certificate = certificateContent
	certificatesModel.PrivateKey = privateKeyContent
	certificatesModel.ExpiredAt = certificateInfo.NotAfter
	certificatesModel.IsEnable = certificateData.IsEnable
	certificatesModel.IsRelease = utils.IsReleaseN
	certificatesModel.Sni = certificateInfo.CommonName

	updateErr := certificatesModel.CertificatesUpdate(id, &certificatesModel)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

type CertificateContent struct {
	ID          string `json:"id"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
	IsEnable    int    `json:"is_enable"`
}

func (c *CertificateContent) CertificateContentInfo(id string) CertificateContent {
	certificateContent := CertificateContent{}

	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	if certificateInfo.ID != id {
		return certificateContent
	}

	certificateContent.ID = certificateInfo.ID
	certificateContent.Certificate = certificateInfo.Certificate
	certificateContent.PrivateKey = certificateInfo.PrivateKey
	certificateContent.IsEnable = certificateInfo.IsEnable

	return certificateContent
}

type CertificateInfo struct {
	ID        string `json:"id"`
	Sni       string `json:"sni"`
	ExpiredAt int64  `json:"expired_at"`
	IsEnable  int    `json:"is_enable"`
	IsRelease int    `json:"is_release"`
}

func (c *CertificateInfo) CertificateListPage(param *validators.CertificateList) ([]CertificateInfo, int, error) {
	certificatesModel := models.Certificates{}
	certificateListInfos, total, certificateListInfosErr := certificatesModel.CertificateListPage(param)

	certificateList := make([]CertificateInfo, 0)
	if len(certificateListInfos) != 0 {
		for _, certificateListInfo := range certificateListInfos {
			certificateInfo := CertificateInfo{}
			certificateInfo.ID = certificateListInfo.ID
			certificateInfo.Sni = certificateListInfo.Sni
			certificateInfo.ExpiredAt = certificateListInfo.ExpiredAt.Unix()
			certificateInfo.IsEnable = certificateListInfo.IsEnable
			certificateInfo.IsRelease = certificateListInfo.IsRelease

			certificateList = append(certificateList, certificateInfo)
		}
	}

	return certificateList, total, certificateListInfosErr
}

func CertificateDelete(id string) error {
	configReleaseErr := CertificateConfigRelease(utils.ReleaseTypeDelete, id)
	if configReleaseErr != nil {
		return configReleaseErr
	}

	certificatesModel := models.Certificates{}
	deleteErr := certificatesModel.CertificateDelete(id)
	if deleteErr != nil {
		CertificateConfigRelease(utils.ReleaseTypePush, id)
		return deleteErr
	}

	return nil
}

func CertificateSwitchEnable(id string, enable int) error {
	certificatesModel := models.Certificates{}
	updateErr := certificatesModel.CertificateSwitchEnable(id, enable)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func CertificateRelease(id string) error {
	certificatesModel := models.Certificates{}
	updateErr := certificatesModel.CertificateSwitchRelease(id, utils.IsReleaseY)
	if updateErr != nil {
		return updateErr
	}

	configReleaseErr := CertificateConfigRelease(utils.ReleaseTypePush, id)
	if configReleaseErr != nil {
		certificatesModel.CertificateSwitchRelease(id, utils.IsReleaseN)
		return configReleaseErr
	}

	return nil
}

func CertificateConfigRelease(releaseType string, certificateId string) error {

	// @todo 获取指定服务的配置数据
	//certificateConfig := generateCertificateConfig(certificateId)

	// @todo 获取数据注册中心对应 服务配置 的key

	fmt.Println("=========certificate release:", releaseType, certificateId)

	// @todo 发布配置到 数据注册中心

	return nil
}

func generateCertificateConfig(serviceId string) string {

	// @todo 根据服务ID 拼接服务的配置数据（主要是用于同步到数据面使用）

	return ""
}
