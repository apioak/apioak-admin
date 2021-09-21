package admin

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
	"strings"
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

func ClusterNodeSwitchEnable(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	var clusterNodeSwitchEnableValidator = validators.ClusterNodeSwitchEnable{}
	if msg, err := packages.ParseRequestParams(c, &clusterNodeSwitchEnableValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	checkClusterNodeNullErr := services.CheckClusterNodeNull(id)
	if checkClusterNodeNullErr != nil {
		utils.Error(c, checkClusterNodeNullErr.Error())
		return
	}

	checkClusterNodeEnableChangeErr := services.CheckClusterNodeEnableChange(id, clusterNodeSwitchEnableValidator.IsEnable)
	if checkClusterNodeEnableChangeErr != nil {
		utils.Error(c, checkClusterNodeEnableChangeErr.Error())
		return
	}

	updateErr := services.ClusterNodeSwitchEnable(id, clusterNodeSwitchEnableValidator.IsEnable)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}
