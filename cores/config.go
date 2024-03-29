package cores

import (
	"apioak-admin/app/packages"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strings"
)

type ConfigServer struct {
	Host string `yaml:"host" mapstructure:"host"`
	Port int    `yaml:"port" mapstructure:"port"`
	Mode string `yaml:"mode" mapstructure:"mode"`
}

type Logger struct {
	LogPath      string `yaml:"log_path" mapstructure:"log_path"`
	LogFileInfo  string `yaml:"log_file_info" mapstructure:"log_file_info"`
	LogFileError string `yaml:"log_file_error" mapstructure:"log_file_error"`
	LogReserve   int64  `yaml:"log_reserve" mapstructure:"log_reserve"`
}

type ConfigDatabase struct {
	Driver             string `yaml:"driver" mapstructure:"driver"`
	Host               string `yaml:"host" mapstructure:"host"`
	Port               int    `yaml:"port" mapstructure:"port"`
	DbName             string `yaml:"db_name" mapstructure:"db_name"`
	Username           string `yaml:"username" mapstructure:"username"`
	Password           string `yaml:"password" mapstructure:"password"`
	MaxIdelConnections int    `yaml:"max_idel_connections" mapstructure:"max_idel_connections"`
	MaxOpenConnections int    `yaml:"max_open_connections" mapstructure:"max_open_connections"`
	SqlMode            bool   `yaml:"sql_mode" mapstructure:"sql_mode"`
}

type ConfigToken struct {
	TokenIssuer string `yaml:"token_issuer" mapstructure:"token_issuer"`
	TokenSecret string `yaml:"token_secret" mapstructure:"token_secret"`
	TokenExpire uint32 `yaml:"token_expire" mapstructure:"token_expire"`
}

type ConfigValidator struct {
	Locale string `yaml:"locale" mapstructure:"locale"`
}

type ConfigApiOak struct {
	Protocol string `yaml:"protocol" mapstructure:"protocol"`
	Ip       string `yaml:"ip" mapstructure:"ip"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Domain   string `yaml:"domain" mapstructure:"domain"`
	Secret   string `yaml:"secret" mapstructure:"secret"`
}

type ConfigRuntime struct {
	DB  *gorm.DB
	Gin *gin.Engine
}

type ConfigGlobal struct {
	Server    ConfigServer    `yaml:"server" mapstructure:"server"`
	Logger    Logger          `yaml:"logger" mapstructure:"logger"`
	Database  ConfigDatabase  `yaml:"database" mapstructure:"database"`
	Validator ConfigValidator `yaml:"validator" mapstructure:"validator"`
	Token     ConfigToken     `yaml:"token"`
	Apioak    ConfigApiOak    `yaml:"apioak" mapstructure:"apioak"`
	Runtime   ConfigRuntime
}

// InitConfig 全局配置初始化
func InitConfig(conf *ConfigGlobal) error {

	filename := ""

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("app")
	v.AddConfigPath("./config/")
	v.SetConfigFile(filename)

	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config is failed err: `%s`", err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(conf); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(conf); err != nil {
		fmt.Println(err)
	}

	protocol := strings.ToLower(conf.Apioak.Protocol)
	if (protocol != "http") && (protocol != "https") {
		protocol = "http"
	}

	packages.SetConfigApiOak(protocol, conf.Apioak.Ip, conf.Apioak.Port, conf.Apioak.Domain, conf.Apioak.Secret)

	return nil
}
