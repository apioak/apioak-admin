package src

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const AppVersion = "0.6.0"

type ConfigApp struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigLog struct {
	Level  string `yaml:"level"`
	Format string `yaml:"json"`
	Error  string `yaml:"error"`
	Access string `yaml:"access"`
}

type ConfigCLI struct {
	Config  string
	Version bool
}

type ConfigEtcd struct {
	Prefix string   `yaml:"prefix"`
	Nodes  []string `yaml:"nodes"`
}

type Config struct {
	CLI  ConfigCLI  `yaml:"-"`
	App  ConfigApp  `yaml:"app"`
	Log  ConfigLog  `yaml:"log"`
	Etcd ConfigEtcd `yaml:"etcd"`
}

var config Config

func init() {
	//flag.StringVar(&config.CLI.Config, "c", "etc/apioak.yaml", "the apioak config file")
	//flag.BoolVar(&config.CLI.Version, "v", false, "the apioak version")
}

func initConfig() error {
	config.CLI.Config = "etc/apioak.yaml"
	configFile, err := ioutil.ReadFile(config.CLI.Config)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return err
	}

	fmt.Println(config.Etcd.Nodes)

	return nil
}
