package middlewares

import (
	"github.com/gin-gonic/gin"
)

func CheckUserLogin(c *gin.Context) {
	// @todo 校验用户是否登录

	c.Next()
}
