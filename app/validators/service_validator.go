package validators

import (
	"encoding/json"
	"strconv"
)

var (
	connectionTimeout    = 3000
	sendTimeout          = 4000
	readTimeout          = 5000
	connectionTimeoutKey = "connection_timeout"
	sendTimeoutKey       = "send_timeout"
	readTimeoutKey       = "read_timeout"
)

type ServiceAdd struct {
	Timeouts     string `json:"timeouts" zh:"超时时间" en:"Time out" binding:"omitempty,json"`
	LoadBalance  int    `json:"load_balance" zh:"负载均衡算法" en:"Load balancing algorithm" binding:"omitempty,oneof=1 2"`
	IsEnable     int    `json:"is_enable" zh:"服务开关" en:"Service enable" binding:"omitempty,oneof=1 2"`
	WebSocket    int    `json:"web_socket" zh:"WebSocket" en:"WebSocket" binding:"omitempty,oneof=1 2"`
	HealthCheck  int    `json:"health_check" zh:"健康检查" en:"Health" binding:"omitempty,oneof=1 2"`
	Protocol     int    `json:"protocol" zh:"请求协议" en:"Protocol" binding:"omitempty,oneof=1 2 3"`
	ServiceNodes string `json:"service_nodes" zh:"上游节点" en:"Service nodes" binding:"required,json,CheckServiceNode"`
	ServiceNames string `json:"service_domains" zh:"域名" en:"Service domains" binding:"required,CheckServiceDomain"`
}

func defaultServiceTimeOut() map[string]uint32 {
	timeInterface := make(map[string]uint32)

	timeInterface[connectionTimeoutKey] = uint32(connectionTimeout)
	timeInterface[sendTimeoutKey] = uint32(sendTimeout)
	timeInterface[readTimeoutKey] = uint32(readTimeout)

	return timeInterface
}

func GetServiceAddTimeOut(times string) string {

	defaultTimeOut := defaultServiceTimeOut()
	if len(times) <= 0 {

		timeStr, err := json.Marshal(defaultTimeOut)
		if err != nil {
			return ""
		}
		return string(timeStr)
	}

	timeInterface := make(map[string]interface{})
	jsonErr := json.Unmarshal([]byte(times), &timeInterface)
	if jsonErr != nil {
		timeStr, err := json.Marshal(defaultTimeOut)
		if err != nil {
			return ""
		}
		return string(timeStr)
	}

	for timeKey, millisecond := range timeInterface {
		switch timeKey {
		case connectionTimeoutKey:

			millisecondInt, err := strconv.Atoi(millisecond.(string))
			if err != nil {
				break
			}
			defaultTimeOut[connectionTimeoutKey] = uint32(millisecondInt)
		case sendTimeoutKey:
			millisecondInt, err := strconv.Atoi(millisecond.(string))
			if err != nil {
				break
			}
			defaultTimeOut[sendTimeoutKey] = uint32(millisecondInt)
		case readTimeoutKey:
			millisecondInt, err := strconv.Atoi(millisecond.(string))
			if err != nil {
				break
			}
			defaultTimeOut[readTimeoutKey] = uint32(millisecondInt)
		}
	}

	timeStr, err := json.Marshal(defaultTimeOut)
	if err != nil {
		return ""
	}
	return string(timeStr)
}
