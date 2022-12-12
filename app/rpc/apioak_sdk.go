package rpc

import (
	"apioak-admin/app/packages"
)

type ApiOak struct {
	Ip     string
	Port   int
	Domain string
	Secret string
}

func NewApiOak() ApiOak {
	apiOak := ApiOak{}
	apiOak.Ip = packages.ConfigApiOak.Ip
	apiOak.Port = packages.ConfigApiOak.Port
	apiOak.Domain = packages.ConfigApiOak.Domain
	apiOak.Secret = packages.ConfigApiOak.Secret

	return apiOak
}

type ConfigObjectName struct {
	Name string `json:"name"`
}

type HealthCheck struct {
	Enabled  bool   `json:"enabled"`
	Tcp      bool   `json:"tcp"`
	Method   string `json:"method"`
	Host     string `json:"host"`
	Uri      string `json:"uri"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type UpstreamNodeConfig struct {
	Name    string      `json:"name"`
	Address string      `json:"address"`
	Port    int         `json:"port"`
	Weight  int         `json:"weight"`
	Health  string      `json:"health"`
	Check   HealthCheck `json:"check"`
}

func (m *ApiOak) UpstreamNodePut(upstreamNodeConfigList []UpstreamNodeConfig) (err error) {
	// @todo 发布逻辑，先检测远程是否存在，存在的直接更新，不存在的直接新增
	return
}

func (m *ApiOak) UpstreamNodeDelete(upstreamNodeConfigList []UpstreamNodeConfig) (err error) {
	// @todo 删除逻辑，先检测远程是否存在，存在的直接删除，不存在忽略
	return
}

type UpstreamConfig struct {
	Name           string             `json:"name"`
	Algorithm      string             `json:"algorithm"`
	ConnectTimeout int                `json:"connect_timeout"`
	WriteTimeout   int                `json:"write_timeout"`
	ReadTimeout    int                `json:"read_timeout"`
	Nodes          []ConfigObjectName `json:"nodes"`
}

func (m *ApiOak) UpstreamPut(upstreamConfigList []UpstreamConfig) (err error) {
	// @todo 发布逻辑，先检测远程是否存在，存在的直接更新，不存在的直接新增
	return
}

func (m *ApiOak) UpstreamDelete(upstreamConfigList []UpstreamConfig) (err error) {
	// @todo 删除逻辑，先检测远程是否存在，存在的直接删除，不存在忽略
	return
}

type RouteConfig struct {
	Name     string             `json:"name"`
	Methods  []string           `json:"methods"`
	Paths    []string           `json:"paths"`
	Enabled  bool               `json:"enabled"`
	Headers  map[string]string  `json:"headers"`
	Service  ConfigObjectName   `json:"service"`
	Upstream ConfigObjectName   `json:"upstream"`
	Plugins  []ConfigObjectName `json:"plugins"`
}

func (m *ApiOak) RoutePut(routeConfigList []RouteConfig) (err error) {
	// @todo 发布逻辑，先检测远程是否存在，存在的直接更新，不存在的直接新增
	return
}

func (m *ApiOak) RouteDelete(routeConfigList []RouteConfig) (err error) {
	// @todo 删除逻辑，先检测远程是否存在，存在的直接删除，不存在忽略
	return
}
