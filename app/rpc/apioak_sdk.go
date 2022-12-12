package rpc

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type ApiOak struct {
	Protocol string
	Ip       string
	Port     int
	Domain   string
	Secret   string
}

var (
	apiOak     *ApiOak
	apiOakOnce sync.Once

	timeOut         = time.Second * 2
	routerUri     = "/apioak/admin/routers"
	upstreamUri     = "/apioak/admin/upstreams"
	upstreamNodeUri = "/apioak/admin/upstream/nodes"
)

func NewApiOak() *ApiOak {

	apiOakOnce.Do(func() {
		apiOak = &ApiOak{
			Protocol: packages.ConfigApiOak.Protocol,
			Ip:       packages.ConfigApiOak.Ip,
			Port:     packages.ConfigApiOak.Port,
			Domain:   packages.ConfigApiOak.Domain,
			Secret:   packages.ConfigApiOak.Secret,
		}
	})

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

func (m *ApiOak) commonPut(resName string, uri string, data interface{}, params url.Values, headers http.Header) (err error) {

	getUri := uri + "/" + resName

	var httpResp utils.HttpResp
	httpResp, err = utils.Get(getUri, params, headers, timeOut)
	if err != nil {
		return
	}

	if httpResp.StatusCode == 404 {
		httpResp, err = utils.PostJson(uri, data, headers, timeOut)
		if err != nil {
			return
		}

		if httpResp.StatusCode != 200 {
			return errors.New(enums.CodeMessages(enums.PublishError))
		}
	} else if httpResp.StatusCode == 200 {
		httpResp, err = utils.PutJson(getUri, data, headers, timeOut)
		if err != nil {
			return
		}

		if httpResp.StatusCode != 200 {
			err = errors.New(enums.CodeMessages(enums.PublishError))
		}
	}

	return
}

func (m *ApiOak) UpstreamNodePut(upstreamNodeConfigList []UpstreamNodeConfig) (err error) {

	if len(upstreamNodeConfigList) == 0 {
		return
	}

	uri := m.Protocol + "://" + m.Ip + ":" + strconv.Itoa(m.Port) + upstreamNodeUri

	for _, upstreamNodeConfigInfo := range upstreamNodeConfigList {

		var param = url.Values{}
		var header = http.Header{}
		if len(m.Domain) > 0 {
			header.Set("Host", m.Domain)
		}

		resName := upstreamNodeConfigInfo.Name
		err = m.commonPut(resName, uri, upstreamNodeConfigInfo, param, header)
		if err != nil {
			return
		}
	}

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
	if len(upstreamConfigList) == 0 {
		return
	}

	uri := m.Protocol + "://" + m.Ip + ":" + strconv.Itoa(m.Port) + upstreamUri

	for _, upstreamConfigInfo := range upstreamConfigList {

		var param = url.Values{}
		var header = http.Header{}
		if len(m.Domain) > 0 {
			header.Set("Host", m.Domain)
		}

		resName := upstreamConfigInfo.Name
		err = m.commonPut(resName, uri, upstreamConfigInfo, param, header)
		if err != nil {
			return
		}
	}

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
	if len(routeConfigList) == 0 {
		return
	}

	uri := m.Protocol + "://" + m.Ip + ":" + strconv.Itoa(m.Port) + routerUri

	for _, routeConfigInfo := range routeConfigList {

		var param = url.Values{}
		var header = http.Header{}
		if len(m.Domain) > 0 {
			header.Set("Host", m.Domain)
		}

		resName := routeConfigInfo.Name
		err = m.commonPut(resName, uri, routeConfigInfo, param, header)
		if err != nil {
			return
		}
	}

	return
}

func (m *ApiOak) RouteDelete(routeConfigList []RouteConfig) (err error) {
	// @todo 删除逻辑，先检测远程是否存在，存在的直接删除，不存在忽略
	return
}
