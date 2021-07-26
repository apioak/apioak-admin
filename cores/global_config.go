package cores

type ConfigGlobal struct {


}

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

type ConfigValidator struct {
	Locale string `yaml:"locale"`
}
