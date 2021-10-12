package middlewares

import (
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"github.com/gin-gonic/gin"
)

func CheckUserLogin(c *gin.Context) {
	token := c.GetHeader("auth-token")

	loginStatus, loginStatusErr := services.CheckUserLoginStatus(token)
	if (loginStatusErr != nil) || (loginStatus == false) {
		utils.Error(c, loginStatusErr.Error())
		c.Abort()
		return
	}

	refresh, refreshErr := services.UserLoginRefresh(token)
	if (refreshErr != nil) || (refresh == false) {
		utils.Error(c, refreshErr.Error())
		c.Abort()
		return
	}

	c.Next()
}
