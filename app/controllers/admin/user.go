package admin

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"github.com/gin-gonic/gin"
)

// UserRegister @todo 用户注册
func UserRegister(c *gin.Context) {

	//获取参数结构体
	var registerStruct = validators.UserRegister{}

	// 参数校验
	if msg, err := packages.ParseRequestParams(c, &registerStruct); err != nil {
		utils.Error(c, msg)
		return
	}

	utils.Ok(c)
}
