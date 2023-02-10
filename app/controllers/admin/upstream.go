package admin

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

func UpstreamList(c *gin.Context) {
	var request = &validators.UpstreamList{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	list, total, err := services.NewServiceUpstream().UpstreamList(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	res := &utils.ResultPage{
		Param:    request,
		Page:     request.Page,
		PageSize: request.PageSize,
		Data:     list,
		Total:    total,
	}

	utils.Ok(c, res)
}

func UpstreamAdd(c *gin.Context) {
	var request = &validators.UpstreamAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	// 初始化默认值
	validators.CorrectUpstreamDefault(request)
	validators.CorrectUpstreamAddNodes(&request.UpstreamNodes)

	serviceUpstream := services.NewServiceUpstream()
	if request.Name != "" {
		err := serviceUpstream.CheckExistName([]string{request.Name}, []string{})
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
	}

	err := serviceUpstream.UpstreamCreate(request)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c)
}
