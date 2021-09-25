package admin

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var userRegisterValidator = validators.UserRegister{}
	if msg, err := packages.ParseRequestParams(c, &userRegisterValidator); err != nil {
		utils.Error(c, msg)
		return
	}

	checkUserEmailExistErr := services.CheckUserEmailExist(userRegisterValidator.Email, []string{})
	if checkUserEmailExistErr != nil {
		utils.Error(c, checkUserEmailExistErr.Error())
		return
	}

	addErr := services.UserCreate(&userRegisterValidator)
	if addErr != nil {
		utils.Error(c, enums.CodeMessages(enums.Error))
		return
	}

	utils.Ok(c)
}

func UserLogin(c *gin.Context) {
	var userLoginValidator = validators.UserLogin{}
	if msg, err := packages.ParseRequestParams(c, &userLoginValidator); err != nil {
		utils.Error(c, msg)
		return
	}



}
