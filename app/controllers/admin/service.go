package admin

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

func ServiceAdd(c *gin.Context) {

	var serviceRegisterStruct = validators.ServiceAdd{}
	if msg, err := packages.ParseRequestParams(c, &serviceRegisterStruct); err != nil {
		utils.Error(c, msg)
		return
	}

	err := services.CheckExistDomain(serviceRegisterStruct.ServiceNames)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	serviceRegisterStruct.Timeouts = validators.GetServiceAddTimeOut(serviceRegisterStruct.Timeouts)
	serviceDomains := validators.GetServiceAddDomains(serviceRegisterStruct.ServiceNames)
	serviceNodes := validators.GetServiceAddNodes(serviceRegisterStruct.ServiceNodes)

	createErr := services.ServiceCreate(&serviceRegisterStruct, &serviceDomains, &serviceNodes)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}
