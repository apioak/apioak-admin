package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
	"strings"
)

func CertificateAdd(c *gin.Context) {
	var certificateAddValidator = validators.CertificateAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &certificateAddValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	addErr := services.CertificateAdd(&certificateAddValidator)
	if addErr != nil {
		utils.Error(c, addErr.Error())
		return
	}

	utils.Ok(c)
}

func CertificateUpdate(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	checkCertificateNull := services.CheckCertificateNull(id)
	if checkCertificateNull != nil {
		utils.Error(c, checkCertificateNull.Error())
		return
	}

	var certificateUpdateValidator = validators.CertificateAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &certificateUpdateValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	updateErr := services.CertificateUpdate(id, &certificateUpdateValidator)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func CertificateInfo(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	checkCertificateNull := services.CheckCertificateNull(id)
	if checkCertificateNull != nil {
		utils.Error(c, checkCertificateNull.Error())
		return
	}

	certificateContent := services.CertificateContent{}
	certificateContentInfo := certificateContent.CertificateContentInfo(id)

	utils.Ok(c, certificateContentInfo)
}

func CertificateList(c *gin.Context) {
	var certificateListValidator = validators.CertificateList{}
	if msg, err := packages.ParseRequestParams(c, &certificateListValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	certificateInfo := services.CertificateInfo{}
	certificateList, total, err := certificateInfo.CertificateListPage(&certificateListValidator)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = certificateListValidator
	result.Page = certificateListValidator.Page
	result.PageSize = certificateListValidator.PageSize
	result.Total = total
	result.Data = certificateList

	utils.Ok(c, result)
}

func CertificateDelete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	checkCertificateNull := services.CheckCertificateNull(id)
	if checkCertificateNull != nil {
		utils.Error(c, checkCertificateNull.Error())
		return
	}

	checkCertificateEnableOnErr := services.CheckCertificateEnableOn(id)
	if checkCertificateEnableOnErr != nil {
		utils.Error(c, checkCertificateEnableOnErr.Error())
		return
	}

	checkCertificateDomainExistErr := services.CheckCertificateDomainExistById(id)
	if checkCertificateDomainExistErr != nil {
		utils.Error(c, checkCertificateDomainExistErr.Error())
		return
	}

	deleteErr := services.CertificateDelete(id)
	if deleteErr != nil {
		utils.Error(c, enums.CodeMessages(enums.Error))
		return
	}

	utils.Ok(c)
}

func CertificateSwitchEnable(c *gin.Context) {

}
