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

func CheckCertificateDelete(id string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)

	if certificateInfo.Release == utils.ReleaseStatusY {
		if certificateInfo.Enable == utils.EnableOn {
			return errors.New(enums.CodeMessages(enums.SwitchONProhibitsOp))
		}
	} else if certificateInfo.Release == utils.ReleaseStatusT {
		return errors.New(enums.CodeMessages(enums.ToReleaseProhibitsOp))
	}

	return nil
}

func CheckCertificateEnableChange(id string, enable int) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	if certificateInfo.Enable == enable {
		return errors.New(enums.CodeMessages(enums.SwitchNoChange))
	}

	return nil
}

func CheckCertificateRelease(id string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	if certificateInfo.Release == utils.ReleaseStatusY {
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
	certificateInfo, certificateInfoErr := utils.DiscernCertificate(&certificateData.Certificate)
	if certificateInfoErr != nil {
		return certificateInfoErr
	}

	checkCertificateExistErr := CheckCertificateExistBySni(certificateInfo.CommonName, "")
	if checkCertificateExistErr != nil {
		return checkCertificateExistErr
	}

	certificatesModel := models.Certificates{
		Certificate: certificateData.Certificate,
		PrivateKey:  certificateData.PrivateKey,
		ExpiredAt:   certificateInfo.NotAfter,
		Enable:      certificateData.IsEnable,
		Release:     utils.ReleaseStatusU,
		Sni:         certificateInfo.CommonName,
	}

	if certificateData.IsRelease == utils.ReleaseY {
		certificatesModel.Release = utils.ReleaseStatusY
	}

	certificateId, addErr := certificatesModel.CertificatesAdd(&certificatesModel)
	if addErr != nil {
		return addErr
	}

	if certificateData.IsRelease == utils.ReleaseY {
		configReleaseErr := CertificateConfigRelease(utils.ReleaseTypePush, certificateId)
		if configReleaseErr != nil {
			certificatesModel.Release = utils.ReleaseStatusU
			certificatesModel.CertificatesUpdate(certificateId, &certificatesModel)
			return configReleaseErr
		}

		return configReleaseErr
	}

	return nil
}

func CertificateUpdate(id string, certificateData *validators.CertificateAddUpdate) error {
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

	certificatesModel.Certificate = certificateData.Certificate
	certificatesModel.PrivateKey = certificateData.PrivateKey
	certificatesModel.ExpiredAt = certificateInfo.NotAfter
	certificatesModel.Enable = certificateData.IsEnable
	certificatesModel.Sni = certificateInfo.CommonName
	if certificateExistInfo.Release == utils.ReleaseStatusY {
		certificatesModel.Release = utils.ReleaseStatusT
	}

	if certificateData.IsRelease == utils.ReleaseY {
		certificatesModel.Release = utils.ReleaseStatusY
	}

	updateErr := certificatesModel.CertificatesUpdate(id, &certificatesModel)
	if updateErr != nil {
		return updateErr
	}

	if certificateData.IsRelease == utils.ReleaseY {
		configReleaseErr := CertificateConfigRelease(utils.ReleaseTypePush, id)
		if configReleaseErr != nil {
			if certificateExistInfo.Release != utils.ReleaseStatusU {
				certificatesModel.Release = utils.ReleaseStatusT
			}
			certificatesModel.CertificatesUpdate(id, &certificatesModel)
			return configReleaseErr
		}

		return configReleaseErr
	}

	return nil
}

type CertificateContent struct {
	ID            string `json:"id"`
	Certificate   string `json:"certificate"`
	PrivateKey    string `json:"private_key"`
	IsEnable      int    `json:"is_enable"`
	ReleaseStatus int    `json:"release_status"`
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
	certificateContent.IsEnable = certificateInfo.Enable
	certificateContent.ReleaseStatus = certificateInfo.Release

	return certificateContent
}

type CertificateInfo struct {
	ID        string `json:"id"`
	Sni       string `json:"sni"`
	ExpiredAt int64  `json:"expired_at"`
	Enable    int    `json:"enable"`
	Release   int    `json:"release"`
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
			certificateInfo.Enable = certificateListInfo.Enable
			certificateInfo.Release = certificateListInfo.Release

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

func CertificateRelease(id string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoById(id)
	updateErr := certificatesModel.CertificateSwitchRelease(id, utils.ReleaseStatusY)
	if updateErr != nil {
		return updateErr
	}

	configReleaseErr := CertificateConfigRelease(utils.ReleaseTypePush, id)
	if configReleaseErr != nil {
		certificatesModel.CertificateSwitchRelease(id, certificateInfo.Release)
		return configReleaseErr
	}

	return nil
}

func CertificateConfigRelease(releaseType string, certificateId string) error {
	certificateConfig, certificateConfigErr := generateCertificateConfig(certificateId)
	if certificateConfigErr != nil {
		return certificateConfigErr
	}

	certificateConfigJson, certificateConfigJsonErr := json.Marshal(certificateConfig)
	if certificateConfigJsonErr != nil {
		return certificateConfigJsonErr
	}
	certificateConfigStr := string(certificateConfigJson)

	etcdKey := utils.EtcdKey(utils.EtcdKeyTypeCertificate, certificateId)
	if len(etcdKey) == 0 {
		return errors.New(enums.CodeMessages(enums.EtcdKeyNull))
	}

	etcdClient := packages.GetEtcdClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*utils.EtcdTimeOut)
	defer cancel()

	var respErr error
	if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Put(ctx, etcdKey, certificateConfigStr)
	} else if strings.ToLower(releaseType) == utils.ReleaseTypePush {
		_, respErr = etcdClient.Delete(ctx, etcdKey)
	}

	if respErr != nil {
		return errors.New(enums.CodeMessages(enums.EtcdUnavailable))
	}

	return nil
}

type CertificateConfig struct {
	ID       string `json:"id"`
	Sni      string `json:"sni"`
	IsEnable int    `json:"is_enable"`
	Pem      string `json:"pem"`
	Key      string `json:"key"`
}

func generateCertificateConfig(id string) (CertificateConfig, error) {
	certificateConfig := CertificateConfig{}
	certificateModel := models.Certificates{}
	certificateInfo := certificateModel.CertificateInfoById(id)
	if len(certificateInfo.ID) == 0 {
		return certificateConfig, errors.New(enums.CodeMessages(enums.CertificateNull))
	}

	certificateConfig.ID = certificateInfo.ID
	certificateConfig.Sni = certificateInfo.Sni
	certificateConfig.IsEnable = certificateInfo.Enable
	certificateConfig.Pem = certificateInfo.Certificate
	certificateConfig.Key = certificateInfo.PrivateKey

	return certificateConfig, nil
}
