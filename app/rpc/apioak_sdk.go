package rpc

import (
	"apioak-admin/app/enums"
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"encoding/json"
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
	Address  string
}

var (
	apiOak     *ApiOak
	apiOakOnce sync.Once

	timeOut         = time.Second * 2
	routerUri       = "/apioak/admin/routers"
	upstreamUri     = "/apioak/admin/upstreams"
	upstreamNodeUri = "/apioak/admin/upstream/nodes"
	serviceUri      = "/apioak/admin/services"
	certificateUri  = "/apioak/admin/certificates"
)

func NewApiOak() *ApiOak {

	apiOakOnce.Do(func() {

		address := packages.ConfigApiOak.Protocol + "://"
		if len(packages.ConfigApiOak.Domain) != 0 {
			address = address + packages.ConfigApiOak.Domain
		} else {
			address = address + packages.ConfigApiOak.Ip
		}
		address = address + ":" + strconv.Itoa(packages.ConfigApiOak.Port)

		apiOak = &ApiOak{
			Protocol: packages.ConfigApiOak.Protocol,
			Ip:       packages.ConfigApiOak.Ip,
			Port:     packages.ConfigApiOak.Port,
			Domain:   packages.ConfigApiOak.Domain,
			Secret:   packages.ConfigApiOak.Secret,
			Address:  address,
		}
	})

	return apiOak
}

type ConfigObjectName struct {
	Id   string `json:"id,omitempty"`
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

	if len(m.Domain) > 0 {
		headers.Set("Host", m.Domain)
	}

	var httpResp utils.HttpResp
	httpResp, err = utils.Get(getUri, params, headers, timeOut)
	if err != nil {
		return
	}

	if httpResp.StatusCode == 404 || httpResp.StatusCode == 500 {
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

	uri := m.Address + upstreamNodeUri

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

func (m *ApiOak) UpstreamNodeDelete(upstreamNodeResIds []string) (err error) {
	if len(upstreamNodeResIds) == 0 {
		return
	}

	uri := m.Address + upstreamNodeUri

	for _, upstreamNodeResId := range upstreamNodeResIds {

		var params = url.Values{}
		var headers = http.Header{}
		if len(m.Domain) > 0 {
			headers.Set("Host", m.Domain)
		}

		uri = uri + "/" + upstreamNodeResId

		var httpResp utils.HttpResp
		httpResp, err = utils.Get(uri, params, headers, timeOut)
		if err != nil {
			return
		}

		if httpResp.StatusCode == 500 {
			err = errors.New(enums.CodeMessages(enums.SyncError))
			return
		} else if httpResp.StatusCode == 200 {
			httpResp, err = utils.Delete(uri, params, headers, timeOut)
			if err != nil {
				return
			}

			if httpResp.StatusCode != 200 {
				err = errors.New(enums.CodeMessages(enums.PublishError))
			}
			return
		}
	}

	return
}

func (m *ApiOak) UpstreamNodeDeleteByIds(upstreamNodeIds []string) (err error) {
	if len(upstreamNodeIds) == 0 {
		return
	}

	uri := m.Address + upstreamNodeUri

	for _, upstreamNodeId := range upstreamNodeIds {

		var params = url.Values{}
		var headers = http.Header{}
		if len(m.Domain) > 0 {
			headers.Set("Host", m.Domain)
		}

		uri = uri + "/" + upstreamNodeId

		var httpResp utils.HttpResp
		httpResp, err = utils.Get(uri, params, headers, timeOut)

		if err != nil {
			return
		}

		if httpResp.StatusCode == 500 {
			err = errors.New(enums.CodeMessages(enums.SyncError))
			return
		} else if httpResp.StatusCode == 200 {
			httpResp, err = utils.Delete(uri, params, headers, timeOut)
			if err != nil {
				return
			}

			if httpResp.StatusCode != 200 {
				err = errors.New(enums.CodeMessages(enums.PublishError))
			}
			return
		}
	}

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

func (m *ApiOak) UpstreamGet(upstreamResIds []string) (list []UpstreamConfig, err error) {
	if len(upstreamResIds) == 0 {
		return
	}

	uri := m.Address + upstreamUri

	for _, upstreamResId := range upstreamResIds {

		var params = url.Values{}
		var headers = http.Header{}
		if len(m.Domain) > 0 {
			headers.Set("Host", m.Domain)
		}

		uri = uri + "/" + upstreamResId

		var httpResp utils.HttpResp
		httpResp, err = utils.Get(uri, params, headers, timeOut)
		if err != nil {
			return
		}

		if httpResp.StatusCode == 404 {
			continue
		} else if httpResp.StatusCode == 500 {
			err = errors.New(enums.CodeMessages(enums.SyncError))
			return
		} else if httpResp.StatusCode == 200 {

			var respData UpstreamConfig
			err = json.Unmarshal(httpResp.Body, &respData)
			if err != nil {
				continue
			}

			if len(respData.Nodes) != 0 {
				for k, node := range respData.Nodes {
					respData.Nodes[k].Name = node.Id
				}
			}

			list = append(list, respData)
		}
	}

	return
}

func (m *ApiOak) UpstreamPut(upstreamConfigList []UpstreamConfig) (err error) {
	if len(upstreamConfigList) == 0 {
		return
	}

	uri := m.Address + upstreamUri

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

func (m *ApiOak) UpstreamDelete(upstreamResIds []string) (err error) {
	if len(upstreamResIds) == 0 {
		return
	}

	uri := m.Address + upstreamUri

	for _, upstreamResId := range upstreamResIds {

		var params = url.Values{}
		var headers = http.Header{}
		if len(m.Domain) > 0 {
			headers.Set("Host", m.Domain)
		}

		uri = uri + "/" + upstreamResId

		var httpResp utils.HttpResp
		httpResp, err = utils.Get(uri, params, headers, timeOut)
		if err != nil {
			return
		}

		if httpResp.StatusCode == 500 {
			err = errors.New(enums.CodeMessages(enums.SyncError))
			return
		} else if httpResp.StatusCode == 200 {
			httpResp, err = utils.Delete(uri, params, headers, timeOut)
			if err != nil {
				return
			}

			if httpResp.StatusCode != 200 {
				err = errors.New(enums.CodeMessages(enums.PublishError))
			}
			return
		}
	}

	return
}

type RouterConfig struct {
	Name     string             `json:"name"`
	Methods  []string           `json:"methods"`
	Paths    []string           `json:"paths"`
	Enabled  bool               `json:"enabled"`
	Headers  map[string]string  `json:"headers"`
	Service  ConfigObjectName   `json:"service"`
	Upstream ConfigObjectName   `json:"upstream"`
	Plugins  []ConfigObjectName `json:"plugins"`
}

func (m *ApiOak) RouterGet(routerResIds []string) (list []RouterConfig, err error) {
	if len(routerResIds) == 0 {
		return
	}

	uri := m.Address + routerUri

	for _, routerResId := range routerResIds {

		var params = url.Values{}
		var headers = http.Header{}
		if len(m.Domain) > 0 {
			headers.Set("Host", m.Domain)
		}

		uri = uri + "/" + routerResId

		var httpResp utils.HttpResp
		httpResp, err = utils.Get(uri, params, headers, timeOut)
		if err != nil {
			return
		}

		if httpResp.StatusCode == 404 {
			continue
		} else if httpResp.StatusCode == 500 {
			err = errors.New(enums.CodeMessages(enums.SyncError))
			return
		} else if httpResp.StatusCode == 200 {

			var respData RouterConfig
			err = json.Unmarshal(httpResp.Body, &respData)
			if err != nil {
				continue
			}

			if len(respData.Upstream.Id) != 0 {
				respData.Upstream.Name = respData.Upstream.Id
			}

			list = append(list, respData)
		}
	}

	return
}

func (m *ApiOak) RouterPut(routerConfigList []RouterConfig) (err error) {
	if len(routerConfigList) == 0 {
		return
	}

	uri := m.Address + routerUri

	for _, routerConfigInfo := range routerConfigList {

		var param = url.Values{}
		var header = http.Header{}
		if len(m.Domain) > 0 {
			header.Set("Host", m.Domain)
		}

		resName := routerConfigInfo.Name
		err = m.commonPut(resName, uri, routerConfigInfo, param, header)
		if err != nil {
			return
		}
	}

	return
}

func (m *ApiOak) RouterDelete(routerResIds []string) (err error) {
	if len(routerResIds) == 0 {
		return
	}

	uri := m.Address + routerUri

	for _, routerResId := range routerResIds {

		var params = url.Values{}
		var headers = http.Header{}
		if len(m.Domain) > 0 {
			headers.Set("Host", m.Domain)
		}

		uri = uri + "/" + routerResId

		var httpResp utils.HttpResp
		httpResp, err = utils.Get(uri, params, headers, timeOut)
		if err != nil {
			return
		}

		if httpResp.StatusCode == 500 {
			err = errors.New(enums.CodeMessages(enums.SyncError))
			return
		} else if httpResp.StatusCode == 200 {
			httpResp, err = utils.Delete(uri, params, headers, timeOut)
			if err != nil {
				return
			}

			if httpResp.StatusCode != 200 {
				err = errors.New(enums.CodeMessages(enums.PublishError))
			}
			return
		}
	}

	return
}

type CertificateInfoResponse struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Sni  []string `json:"sni"`
	Cert string   `json:"cert"`
	Key  string   `json:"key"`
}

func (m *ApiOak) CertificateInfo(resID string) (CertificateInfoResponse, error) {

	var params = url.Values{}
	var headers = http.Header{}
	if len(m.Domain) > 0 {
		headers.Set("Host", m.Domain)
	}

	uri := m.Address + certificateUri + "/" + resID

	httpResp, err := utils.Get(uri, params, headers, timeOut)
	if err != nil || httpResp.StatusCode != 200 {
		packages.Log.Error("Failed to obtain the data side certificate information", err)
		return CertificateInfoResponse{}, err
	}

	var body CertificateInfoResponse
	err = json.Unmarshal(httpResp.Body, &body)

	if err != nil {
		packages.Log.Error("Failed to parse data side certificate information", err)
		return CertificateInfoResponse{}, err
	}

	return body, nil

}

func (m *ApiOak) CertificateDelete(resID string) error {
	var params = url.Values{}
	var headers = http.Header{}

	uri := m.Address + certificateUri + "/" + resID

	httpResp, err := utils.Get(uri, params, headers, timeOut)
	if err != nil {
		packages.Log.Error("[delete]:Failed to obtain the data side certificate information", err)
		return errors.New(err.Error())
	}

	if httpResp.StatusCode != 200 {
		return nil
	}

	dHttpResp, err := utils.Delete(uri, params, headers, timeOut)

	if err != nil || dHttpResp.StatusCode != 200 {
		packages.Log.Error("[delete]:Failed to delete the data side certificate information", err)
		return errors.New(err.Error())
	}

	return nil
}

type CertificateReleaseRequest struct {
	Name string   `json:"name"`
	Sni  []string `json:"snis"`
	Cert string   `json:"cert"`
	Key  string   `json:"key"`
}

func (m *ApiOak) CertificateRelease(request *CertificateReleaseRequest) error {

	resName := request.Name
	uri := m.Address + certificateUri
	err := m.commonPut(resName, uri, request, url.Values{}, http.Header{})

	if err != nil {
		return err
	}

	return nil
}
