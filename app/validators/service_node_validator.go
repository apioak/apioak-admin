package validators

import (
	"apioak-admin/app/packages"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	nodeLocalEn            = "en"
	nodeLocalZh            = "zh"
	nodeLocalErrorMessages = map[string]map[string]string{
		nodeLocalEn: {
			"required": " is a required field",
			"ip":       " must be a valid IP address",
			"max":      " must be %d or less",
			"min":      " must be %d or greater",
		},
		nodeLocalZh: {
			"required": "为必填字段",
			"ip":       "必须是一个有效的IP地址",
			"max":      "必须小于或等于%d",
			"min":      "最小只能为%d",
		},
	}
	nodePortMin   = 0
	nodePortMax   = 65535
	nodeWeightMin = 0
	nodeWeightMax = 100
	serviceNodes  = make([]ServiceNodeAdd, 0)
)

type ServiceNodeAdd struct {
	NodeIp     string `json:"node_ip" zh:"上游节点IP" en:"Node IP" binding:"required,ip"`
	NodePort   int    `json:"node_port" zh:"端口" en:"Node port" binding:"omitempty,min=0,max=10000000000000"`
	NodeWeight int    `json:"node_weight" zh:"权重" en:"Node weight" binding:"omitempty,min=0,max=100"`
}

func CheckServiceNode(fl validator.FieldLevel) bool {

	serviceNodeIps := fl.Field().String()
	var nodeIps []interface{}

	jsonErr := json.Unmarshal([]byte(serviceNodeIps), &nodeIps)
	if jsonErr != nil {
		packages.SetAllCustomizeValidatorErrMsgs("CheckServiceNode", ("json.Unmarshal error: [" + jsonErr.Error() + "]"))
		return false
	}

	if len(nodeIps) <= 0 {
		packages.SetAllCustomizeValidatorErrMsgs("CheckServiceNode", ("json.Unmarshal: parsed ip is empty"))
		return false
	}

	customizeValidator := packages.GetCustomizeValidator()
	for _, ipInfos := range nodeIps {
		switch ipInfo := ipInfos.(type) {
		case map[string]interface{}:

			nodeIp := ""
			nodePort := 80
			nodeWeight := 100

			if ipInfo["node_ip"] != nil {
				nodeIp = ipInfo["node_ip"].(string)
			}
			if ipInfo["node_port"] != nil {
				nodePort = int(ipInfo["node_port"].(float64))
			}
			if ipInfo["node_weight"] != nil {
				nodeWeight = int(ipInfo["node_weight"].(float64))
			}

			nodeIP := ServiceNodeAdd{
				NodeIp:     nodeIp,
				NodePort:   nodePort,
				NodeWeight: nodeWeight,
			}

			nodeIPErr := customizeValidator.Struct(nodeIP)
			if nodeIPErr != nil {
				var (
					structField string
					tag         string
					field       string
					errMsg      string
				)

				for _, e := range nodeIPErr.(validator.ValidationErrors) {
					structField = e.StructField()
					tag = e.Tag()
					field = e.Field()
					break
				}

				switch strings.ToLower(structField) {
				case "nodeip":
					errMsg = nodeIpValidator(tag, field)
				case "nodeport":
					errMsg = nodePortValidator(tag, field)
				case "nodeweight":
					errMsg = nodeWeightValidator(tag, field)
				}
				packages.SetAllCustomizeValidatorErrMsgs("CheckServiceNode", errMsg)
				serviceNodes = nil
				return false
			}

			serviceNodes = append(serviceNodes, nodeIP)
		default:
			packages.SetAllCustomizeValidatorErrMsgs("CheckServiceNode", ("json.Unmarshal: parsed ip type is wrong"))
			serviceNodes = nil
			return false
		}
	}

	return true
}

func nodeIpValidator(tag string, field string) string {
	return field + nodeLocalErrorMessages[strings.ToLower(packages.GetValidatorLocale())][strings.ToLower(tag)]
}

func nodePortValidator(tag string, field string) string {
	var errMsg string

	switch strings.ToLower(tag) {
	case "min":
		errMsg = fmt.Sprintf(field+nodeLocalErrorMessages[strings.ToLower(packages.GetValidatorLocale())][strings.ToLower(tag)], nodePortMin)
	case "max":
		errMsg = fmt.Sprintf(field+nodeLocalErrorMessages[strings.ToLower(packages.GetValidatorLocale())][strings.ToLower(tag)], nodePortMax)
	}
	return errMsg
}

func nodeWeightValidator(tag string, field string) string {
	var errMsg string

	switch strings.ToLower(tag) {
	case "min":
		errMsg = fmt.Sprintf(field+nodeLocalErrorMessages[strings.ToLower(packages.GetValidatorLocale())][strings.ToLower(tag)], nodeWeightMin)
	case "max":
		errMsg = fmt.Sprintf(field+nodeLocalErrorMessages[strings.ToLower(packages.GetValidatorLocale())][strings.ToLower(tag)], nodeWeightMax)
	}
	return errMsg
}

func GetServiceNodeInfo() []ServiceNodeAdd {
	return serviceNodes
}
