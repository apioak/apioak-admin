package src

import (
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Application struct {
		Listen int  `yaml:"listen"`
		Debug  bool `yaml:"debug"`
	} `yaml:"application"`
	Etcd struct {
		Prefix string   `yaml:"prefix"`
		Nodes  []string `yaml:"nodes"`
	} `yaml:"etcd"`
	Mysql struct {
		Host     string `yaml:"host"`
		Database string `yaml:"database"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password int    `yaml:"password"`
		Prefix   string `yaml:"prefix"`
	} `yaml:"mysql"`
}

var config Config

func initConfig() error {
	var app = cli.NewApp()
	app.Name = "apioak-admin"
	app.Usage = "Administrator control panel"
	app.Version = AppVersion

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "conf/apioak-admin.yaml",
			Usage: "set configuration file",
		},
	}

	app.Action = func(c *cli.Context) error {
		_, err := os.Stat(c.String("config"))
		if err != nil {
			return err
		}

		configFile, err := ioutil.ReadFile(c.String("config"))

		if err != nil {
			return err
		}

		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}

	return nil
}
