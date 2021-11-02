package middlewares

import (
	"apioak-admin/app/services"
	"apioak-admin/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckUserLogin(c *gin.Context) {
	token := c.GetHeader("auth-token")

	loginStatus, loginStatusErr := services.CheckUserLoginStatus(token)
	if (loginStatusErr != nil) || (loginStatus == false) {
		utils.CustomError(c, http.StatusUnauthorized, loginStatusErr.Error())
		c.Abort()
		return
	}

	refresh, refreshErr := services.UserLoginRefresh(token)
	if (refreshErr != nil) || (refresh == false) {
		utils.CustomError(c, http.StatusUnauthorized, refreshErr.Error())
		c.Abort()
		return
	}

	c.Next()
}
