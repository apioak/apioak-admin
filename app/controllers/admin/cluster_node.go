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

func ClusterNodeAdd(c *gin.Context) {
	var bindParams = validators.ClusterNodeAdd{
		NodeStatus: utils.ClusterNodeStatusUnhealthy,
	}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	checkClusterNodeExistErr := services.CheckClusterNodeExist(bindParams.NodeIP)
	if checkClusterNodeExistErr != nil {
		utils.Error(c, checkClusterNodeExistErr.Error())
		return
	}

	addErr := services.ClusterNodeAdd(&bindParams)
	if addErr != nil {
		utils.Error(c, addErr.Error())
		return
	}

	utils.Ok(c)
}

func ClusterNodeList(c *gin.Context) {
	var bindParams = validators.ClusterNodeList{}
	if msg, err := packages.ParseRequestParams(c, &bindParams); err != nil {
		utils.Error(c, msg)
		return
	}

	clusterNodeListInfo := services.ClusterNodeListInfo{}
	clusterNodeList, total, err := clusterNodeListInfo.ClusterNodeListPage(&bindParams)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	result := utils.ResultPage{}
	result.Param = bindParams
	result.Page = bindParams.Page
	result.PageSize = bindParams.PageSize
	result.Total = total
	result.Data = clusterNodeList

	utils.Ok(c, result)
}

func ClusterNodeDelete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	checkClusterNodeNullErr := services.CheckClusterNodeNull(id)
	if checkClusterNodeNullErr != nil {
		utils.Error(c, checkClusterNodeNullErr.Error())
		return
	}

	deleteErr := services.ClusterNodeDelete(id)
	if deleteErr != nil {
		utils.Error(c, enums.CodeMessages(enums.Error))
		return
	}

	utils.Ok(c)
}
