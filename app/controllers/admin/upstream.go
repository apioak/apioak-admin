package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
	"strings"
)

func UpstreamList(c *gin.Context) {
	var request = &validators.UpstreamList{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	list, total, err := services.NewServiceUpstream().UpstreamListPage(request)
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

func UpstreamNameList(c *gin.Context) {
	upstreamModel := models.Upstreams{}
	upstreamNameList, err := upstreamModel.UpstreamReleaseNameList()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, upstreamNameList)
}

func UpstreamInfo(c *gin.Context) {
	resId := strings.TrimSpace(c.Param("res_id"))

	if resId == "" {
		utils.Error(c, enums.CodeMessages(enums.ParamsError))
		return
	}

	upstreamInfo, err := services.NewServiceUpstream().UpstreamInfoByResId(resId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Ok(c, upstreamInfo)
}
