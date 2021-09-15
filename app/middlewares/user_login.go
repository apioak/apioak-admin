package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CheckUserLogin(c *gin.Context) {
	// @todo 校验用户是否登录
	token := c.GetHeader("token")

	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~ check-login-middleware:[" + token + "] ~~~~~~~~~~~~~~~~~~~~~~~~~~")

	c.Next()
}
