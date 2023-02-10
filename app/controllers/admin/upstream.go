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

func UpstreamUpdate(c *gin.Context) {
	var request = &validators.UpstreamAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	validators.CorrectUpstreamDefault(request)
	validators.CorrectUpstreamAddNodes(&request.UpstreamNodes)

	resId := strings.TrimSpace(c.Param("res_id"))

	serviceUpstream := services.NewServiceUpstream()
	checkUpstreamExistErr := serviceUpstream.CheckUpstreamExist(resId)
	if checkUpstreamExistErr != nil {
		utils.Error(c, checkUpstreamExistErr.Error())
		return
	}

	if request.Name != "" {
		err := serviceUpstream.CheckExistName([]string{request.Name}, []string{resId})
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
	}

	updateErr := serviceUpstream.UpstreamUpdate(resId, request)
	if updateErr != nil {
		utils.Error(c, updateErr.Error())
		return
	}

	utils.Ok(c)
}

func UpstreamDelete(c *gin.Context) {
	resId := strings.TrimSpace(c.Param("res_id"))

	serviceUpstream := services.NewServiceUpstream()
	checkUpstreamExistErr := serviceUpstream.CheckUpstreamExist(resId)
	if checkUpstreamExistErr != nil {
		utils.Error(c, checkUpstreamExistErr.Error())
		return
	}

	checkUpstreamUseErr := serviceUpstream.CheckUpstreamUse(resId)
	if checkUpstreamUseErr != nil {
		utils.Error(c, checkUpstreamUseErr.Error())
		return
	}

	deleteErr := serviceUpstream.UpstreamDelete(resId)
	if deleteErr != nil {
		utils.Error(c, deleteErr.Error())
		return
	}

	utils.Ok(c)
}

func UpstreamUpdateName(c *gin.Context) {
	var request = &validators.UpstreamUpdateName{}
	if msg, err := packages.ParseRequestParams(c, request); err != nil {
		utils.Error(c, msg)
		return
	}

	resId := strings.TrimSpace(c.Param("res_id"))

	serviceUpstream := services.NewServiceUpstream()
	checkUpstreamExistErr := serviceUpstream.CheckUpstreamExist(resId)
	if checkUpstreamExistErr != nil {
		utils.Error(c, checkUpstreamExistErr.Error())
		return
	}

	if request.Name != "" {
		err := serviceUpstream.CheckExistName([]string{request.Name}, []string{resId})
		if err != nil {
			utils.Error(c, err.Error())
			return
		}
	}

	upstreamModel := models.Upstreams{}
	updateNameErr := upstreamModel.UpstreamUpdateName(resId, request.Name)
	if updateNameErr != nil {
		utils.Error(c, updateNameErr.Error())
		return
	}

	utils.Ok(c)
}
