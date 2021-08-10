package admin

import (
	"apioak-admin/app/packages"
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

	utils.Ok(c, serviceRegisterStruct)
}
