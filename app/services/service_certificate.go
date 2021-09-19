package services

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"encoding/json"
	"errors"
)

func CheckCertificateExist(sni string, filterId string) error {
	certificatesModel := models.Certificates{}
	certificateInfo := certificatesModel.CertificateInfoBySni(sni, filterId)
	if certificateInfo.Sni == sni {
		return errors.New(enums.CodeMessages(enums.CertificateExist))
	}

	return nil
}

func CertificateAdd(certificateData *validators.CertificateAdd) error {
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

	checkCertificateExistErr := CheckCertificateExist(certificateInfo.CommonName, "")
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
