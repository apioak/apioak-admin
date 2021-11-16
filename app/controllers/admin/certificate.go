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
	var bindParams = validators.CertificateAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	addErr := services.CertificateAdd(&bindParams)
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

	var bindParams = validators.CertificateAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	updateErr := services.CertificateUpdate(id, &bindParams)
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
	var bindParams = validators.CertificateList{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	certificateInfo := services.CertificateInfo{}
	certificateList, total, err := certificateInfo.CertificateListPage(&bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = bindParams
	result.Page = bindParams.Page
	result.PageSize = bindParams.PageSize
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

	checkCertificateDeleteErr := services.CheckCertificateDelete(id)
	if checkCertificateDeleteErr != nil {
		utils.Error(c, checkCertificateDeleteErr.Error())
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
	id := strings.TrimSpace(c.Param("id"))

	var bindParams = validators.CertificateSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkCertificateNullErr := services.CheckCertificateNull(id)
	if checkCertificateNullErr != nil {
		utils.Error(c, checkCertificateNullErr.Error())
		return
	}

	checkCertificateEnableChangeErr := services.CheckCertificateEnableChange(id, bindParams.IsEnable)
	if checkCertificateEnableChangeErr != nil {
		utils.Error(c, checkCertificateEnableChangeErr.Error())
		return
	}

	if bindParams.IsEnable == utils.EnableOff {
		checkCertificateDomainExistErr := services.CheckCertificateDomainExistById(id)
		if checkCertificateDomainExistErr != nil {
			utils.Error(c, checkCertificateDomainExistErr.Error())
			return
		}
	}

	updateErr := services.CertificateSwitchEnable(id, bindParams.IsEnable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func CertificateSwitchRelease(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	checkCertificateNullErr := services.CheckCertificateNull(id)
	if checkCertificateNullErr != nil {
		utils.Error(c, checkCertificateNullErr.Error())
		return
	}

	checkCertificateReleaseErr := services.CheckCertificateRelease(id)
	if checkCertificateReleaseErr != nil {
		utils.Error(c, checkCertificateReleaseErr.Error())
		return
	}

	certificateReleaseErr := services.CertificateRelease(id)
	if certificateReleaseErr != nil {
		utils.Error(c, certificateReleaseErr.Error())
		return
	}

	utils.Ok(c)
}
