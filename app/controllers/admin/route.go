package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/models"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

func RouteAdd(c *gin.Context) {

	var routeAddUpdateValidator = validators.RouteAddUpdate{}
	if msg, err := packages.ParseRequestParams(c, &routeAddUpdateValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	if routeAddUpdateValidator.RoutePath == utils.DefaultRoutePath {
		utils.Error(c, enums.CodeMessages(enums.RouteDefaultPathNoPermission))
		return
	}

	err := services.CheckExistRoutePath(routeAddUpdateValidator.RoutePath, []string{})
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	serviceModel := &models.Services{}
	serviceInfo := serviceModel.ServiceInfoById(routeAddUpdateValidator.ServiceID)
	if len(serviceInfo.ID) == 0 {
		utils.Error(c, enums.CodeMessages(enums.ServiceNull))
		return
	}

	validators.GetRouteAttributesDefault(&routeAddUpdateValidator)

	createErr := services.RouteCreate(&routeAddUpdateValidator)
	if createErr != nil {
		utils.Error(c, createErr.Error())
		return
	}

	utils.Ok(c)
}
