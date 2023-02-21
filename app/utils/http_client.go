package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpResp struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func PostJson(uri string, params interface{}, header http.Header, timeout time.Duration) (httpResp HttpResp, err error) {
	postParam, err := json.Marshal(params)
	if err != nil {
		return
	}
	header.Add("Content-Type", "application/json")

	return HttpDo("POST", uri, string(postParam), header, timeout)
}

func PutJson(uri string, params interface{}, header http.Header, timeout time.Duration) (httpResp HttpResp, err error) {
	postParam, err := json.Marshal(params)
	if err != nil {
		return
	}
	header.Add("Content-Type", "application/json")

	return HttpDo("PUT", uri, string(postParam), header, timeout)
}

func PostForm(uri string, params url.Values, header http.Header, timeout time.Duration) (httpResp HttpResp, err error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")

	return HttpDo("POST", uri, params.Encode(), header, timeout)
}

func Get(uri string, params url.Values, header http.Header, timeout time.Duration) (httpResp HttpResp, err error) {
	if len(params) > 0 {
		return HttpDo("GET", uri+"?"+params.Encode(), "", header, timeout)
	}

	return HttpDo("GET", uri, "", header, timeout)
}

func Delete(uri string, params url.Values, header http.Header, timeout time.Duration) (httpResp HttpResp, err error) {
	if len(params) > 0 {
		return HttpDo("DELETE", uri+"?"+params.Encode(), "", header, timeout)
	}

	return HttpDo("DELETE", uri, "", header, timeout)
}

func HttpDo(method, uri, params string, header http.Header, timeout time.Duration) (httpResp HttpResp, err error) {
	client := &http.Client{}
	client.Timeout = timeout

	request, err := http.NewRequest(method, uri, strings.NewReader(params))
	if err != nil {
		return
	}

	if host, ok := header["Host"]; ok {
		request.Host = host[0]
		delete(header, "Host")
	}
	request.Header = header
	request.Close = true

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	httpResp.StatusCode = resp.StatusCode
	httpResp.Header = resp.Header

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	httpResp.Body = body

	return
}
