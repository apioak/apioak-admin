package admin

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

func ClusterNodeList(c *gin.Context) {
	var clusterNodeListValidator = validators.ClusterNodeList{}
	if msg, err := packages.ParseRequestParams(c, &clusterNodeListValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	clusterNodeListInfo := services.ClusterNodeListInfo{}
	clusterNodeList, total, err := clusterNodeListInfo.ClusterNodeListPage(&clusterNodeListValidator)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = clusterNodeListValidator
	result.Page = clusterNodeListValidator.Page
	result.PageSize = clusterNodeListValidator.PageSize
	result.Total = total
	result.Data = clusterNodeList

	utils.Ok(c, result)

}
