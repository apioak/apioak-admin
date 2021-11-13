package utils

import (
	"apioak-admin/app/enums"
	"github.com/gin-gonic/gin"
	"net/http"
)

type result struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 状态码信息
	Data interface{} `json:"data"` // 结果数据
}

type ResultPage struct {
	Param    interface{} `json:"param"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int         `json:"total"`
	Data     interface{} `json:"data"`
}

func Ok(c *gin.Context, data ...interface{}) {
	resultMsg := &result{}
	resultMsg.Code = enums.Success
	resultMsg.Msg = enums.CodeMessages(enums.Success)
	if len(data) > 0 {
		resultMsg.Data = data[0]
	}
	Response(c, resultMsg)
}

func Error(c *gin.Context, message string) {
	resultMsg := &result{}
	resultMsg.Code = enums.Error
	resultMsg.Msg = message
	Response(c, resultMsg)
}

func CustomError(c *gin.Context, code int, message string) {
	resultMsg := &result{}
	resultMsg.Code = code
	resultMsg.Msg = message
	Response(c, resultMsg)
}

func Response(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, result)
}
