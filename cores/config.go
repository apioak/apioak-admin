package cores

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io/ioutil"
)

type ConfigServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type ConfigDatabase struct {
	Driver             string `yaml:"driver"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	DbName             string `yaml:"db_name"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	MaxIdelConnections int    `yaml:"max_idel_connections"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
	SqlMode            bool   `yaml:"sql_mode"`
}

type ConfigEtcd struct {
	HostPort string `yaml:"host_port"`
}

type ConfigToken struct {
	TokenIssuer string `yaml:"token_issuer"`
	TokenSecret string `yaml:"token_secret"`
	TokenExpire uint32 `yaml:"token_expire"`
}

type ConfigValidator struct {
	Locale string `yaml:"locale"`
}

type ConfigRuntime struct {
	DB  *gorm.DB
	Gin *gin.Engine
}

type ConfigGlobal struct {
	Server    ConfigServer    `yaml:"server"`
	Database  ConfigDatabase  `yaml:"database"`
	Etcd      ConfigEtcd      `yaml:"etcd"`
	Validator ConfigValidator `yaml:"validator"`
	Token     ConfigToken     `yaml:"token"`
	Runtime   ConfigRuntime
}

// InitConfig 全局配置初始化
func InitConfig(conf *ConfigGlobal) error {

	// 读取配置文件
	content, err := ioutil.ReadFile("config/app.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, conf)
	if err != nil {
		return err
	}

	return nil
}
