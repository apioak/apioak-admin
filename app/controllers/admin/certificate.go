package admin

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

func CertificateAdd(c *gin.Context) {
	var certificateAddValidator = validators.CertificateAdd{}
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
