package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"encoding/json"
	"errors"
	"strings"
	"time"
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

func CheckCertificateEnableOn(id string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	if certificateInfo.IsEnable == utils.EnableOn {
		return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
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
	ID        string    `json:"id"`
	Sni       string    `json:"sni"`
	ExpiredAt time.Time `json:"expired_at"`
	IsEnable  int       `json:"is_enable"`
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
			certificateInfo.ExpiredAt = certificateListInfo.ExpiredAt
			certificateInfo.IsEnable = certificateListInfo.IsEnable

			certificateList = append(certificateList, certificateInfo)
		}
	}

	return certificateList, total, certificateListInfosErr
}

func CertificateDelete(id string) error {
	certificatesModel := models.Certificates{}
	deleteErr := certificatesModel.CertificateDelete(id)
	if deleteErr != nil {
		return deleteErr
	}

	// @todo 需要同步远程数据中心

	return nil
}

func CertificateSwitchEnable(id string, enable int) error {
	certificatesModel := models.Certificates{}
	updateErr := certificatesModel.CertificateSwitchEnable(id, enable)
	if updateErr != nil {
		return updateErr
	}

	// @todo 触发远程发布数据

	return nil
}
