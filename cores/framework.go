package cores

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func InitFramework(conf *ConfigGlobal) error {
	switch strings.ToLower(conf.Server.Mode) {
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	conf.Runtime.Gin = gin.Default()
	return nil
}

func RunFramework(conf *ConfigGlobal) error {
	return conf.Runtime.Gin.Run(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
}
