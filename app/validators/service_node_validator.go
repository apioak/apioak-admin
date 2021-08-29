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
	defaultNodePort   = 80
	defaultNodeWeight = 100
	nodePortMin       = 0
	nodePortMax       = 65535
	nodeWeightMin     = 0
	nodeWeightMax     = 100
)

type ServiceNodeAddUpdate struct {
	NodeIp     string `json:"node_ip" zh:"上游节点IP" en:"Node IP" binding:"required,ip"`
	NodePort   int    `json:"node_port" zh:"端口" en:"Node port" binding:"omitempty,min=0,max=65535"`
	NodeWeight int    `json:"node_weight" zh:"权重" en:"Node weight" binding:"omitempty,min=0,max=100"`
}

func parsingNodeIpInfos(nodeInfosString string) ([]interface{}, error) {
	var nodeInfos []interface{}
	jsonErr := json.Unmarshal([]byte(nodeInfosString), &nodeInfos)
	if jsonErr != nil {
		return nil, fmt.Errorf("json.Unmarshal error: [" + jsonErr.Error() + "]")
	}

	if len(nodeInfos) <= 0 {
		return nil, fmt.Errorf("json.Unmarshal: parsed ip is empty")
	}

	return nodeInfos, nil
}

func CheckServiceNode(fl validator.FieldLevel) bool {
	serviceNodeIpsString := fl.Field().String()
	nodeIps, err := parsingNodeIpInfos(serviceNodeIpsString)
	if err != nil {
		packages.SetAllCustomizeValidatorErrMsgs("CheckServiceNode", err.Error())
		return false
	}

	customizeValidator := packages.GetCustomizeValidator()
	for _, ipInfos := range nodeIps {
		nodeIpPortWeight, ok := ipInfos.(map[string]interface{})
		if ok == false {
			packages.SetAllCustomizeValidatorErrMsgs("CheckServiceNode", ("json.Unmarshal: parsed ip type is wrong"))
			return false
		}

		nodeIp := ""
		nodePort := defaultNodePort
		nodeWeight := defaultNodeWeight

		if nodeIpPortWeight["node_ip"] != nil {
			nodeIp = nodeIpPortWeight["node_ip"].(string)
		}
		if nodeIpPortWeight["node_port"] != nil {
			nodePort = int(nodeIpPortWeight["node_port"].(float64))
		}
		if nodeIpPortWeight["node_weight"] != nil {
			nodeWeight = int(nodeIpPortWeight["node_weight"].(float64))
		}

		nodeIP := ServiceNodeAddUpdate{
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

func GetServiceAddNodes(serviceNodesString string) []ServiceNodeAddUpdate {
	serviceNodes := []ServiceNodeAddUpdate{}
	nodeIpInfos, err := parsingNodeIpInfos(serviceNodesString)
	if err != nil {
		return serviceNodes
	}

	for _, nodeIpInfo := range nodeIpInfos {
		nodeIpPortWeight, ok := nodeIpInfo.(map[string]interface{})
		if ok == false {
			continue
		}

		nodeIp := ""
		nodePort := defaultNodePort
		nodeWeight := defaultNodeWeight

		if nodeIpPortWeight["node_ip"] != nil {
			nodeIp = nodeIpPortWeight["node_ip"].(string)
		}
		if nodeIpPortWeight["node_port"] != nil {
			nodePort = int(nodeIpPortWeight["node_port"].(float64))
		}
		if nodeIpPortWeight["node_weight"] != nil {
			nodeWeight = int(nodeIpPortWeight["node_weight"].(float64))
		}

		nodeIPPortWeight := ServiceNodeAddUpdate{
			NodeIp:     nodeIp,
			NodePort:   nodePort,
			NodeWeight: nodeWeight,
		}
		serviceNodes = append(serviceNodes, nodeIPPortWeight)
	}

	return serviceNodes
}
