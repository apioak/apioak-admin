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

	var request = &validators.CertificateAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}
	err := services.NewCertificateService().CertificateAdd(request)

	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func CertificateUpdate(c *gin.Context) {
	resID := strings.TrimSpace(c.Param("id"))

	var bindParams = &validators.CertificateAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.NewCertificateService().CertificateUpdate(resID, bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}

func CertificateInfo(c *gin.Context) {
	resID := strings.TrimSpace(c.Param("id"))

	if resID == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	res, err := services.NewCertificateService().CertificateInfo(resID)

	if err != nil {
		utils.Error(c, err.Error())
	}

	utils.Ok(c, res)
}

func CertificateList(c *gin.Context) {
	var bindParams = validators.CertificateList{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	list, total, err := services.NewCertificateService().CertificateListPage(&bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, utils.ResultPage{
		Param:    bindParams,
		Page:     bindParams.Page,
		PageSize: bindParams.PageSize,
		Total:    total,
		Data:     list,
	})
}

func CertificateDelete(c *gin.Context) {
	resID := strings.TrimSpace(c.Param("id"))

	err := services.NewCertificateService().CertificateDelete(resID)
	if err != nil {
		utils.Error(c, enums.CodeMessages(enums.Error))
		return
	}

	utils.Ok(c)
}

func CertificateSwitchEnable(c *gin.Context) {
	resID := strings.TrimSpace(c.Param("id"))

	var bindParams = validators.CertificateSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.NewCertificateService().CertificateSwitchEnable(resID, bindParams.Enable)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}
