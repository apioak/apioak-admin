package validators

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	nodeLocalErrorMessages = map[string]map[string]string{
		utils.LocalEn: {
			"required": " is a required field",
			"ip":       " must be a valid IP address",
			"max":      " must be %d or less",
			"min":      " must be %d or greater",
		},
		utils.LocalZh: {
			"required": "为必填字段",
			"ip":       "必须是一个有效的IP地址",
			"max":      "必须小于或等于%d",
			"min":      "最小只能为%d",
		},
	}
	defaultNodePort   = 80
	nodePortMin       = 0
	nodePortMax       = 65535
	nodeWeightMin     = 0
	nodeWeightMax     = 100
)

type ServiceNodeAddUpdate struct {
	NodeIp     string `json:"node_ip" zh:"上游节点IP" en:"Node IP" binding:"required,ip"`
	NodePort   int    `json:"node_port" zh:"端口" en:"Node port" binding:"omitempty,min=1,max=65535"`
	NodeWeight int    `json:"node_weight" zh:"权重" en:"Node weight" binding:"omitempty,min=0,max=100"`
}

func CheckServiceNode(fl validator.FieldLevel) bool {
	serviceNodeIpsInterface := fl.Field().Interface()
	serviceNodeIps := serviceNodeIpsInterface.([]ServiceNodeAddUpdate)

	customizeValidator := packages.GetCustomizeValidator()
	for _, serviceNodeIpInfo := range serviceNodeIps {

		serviceNodeIpInfo.NodeIp = strings.TrimSpace(serviceNodeIpInfo.NodeIp)
		if serviceNodeIpInfo.NodePort == 0 {
			serviceNodeIpInfo.NodePort = defaultNodePort
		}

		nodeIPErr := customizeValidator.Struct(serviceNodeIpInfo)
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

			switch strings.ToUpper(structField) {
			case "NODEIP":
				errMsg = nodeIpValidator(tag, field)
			case "NODEPORT":
				errMsg = nodePortValidator(tag, field)
			case "NODEWEIGHT":
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

func CorrectServiceAddNodes(serviceNodes *[]ServiceNodeAddUpdate) {
	tmpNodeIpMap := make(map[string]byte, 0)
	tmpNodeInfos := make([]ServiceNodeAddUpdate, 0)

	for _, nodeIpInfo := range *serviceNodes {
		nodeIpTrim := strings.TrimSpace(nodeIpInfo.NodeIp)
		if len(nodeIpTrim) <= 0 {
			continue
		}

		_, exist := tmpNodeIpMap[nodeIpTrim]
		if exist {
			continue
		}

		tmpNodeIpMap[nodeIpTrim] = 0
		tmpNodeInfos = append(tmpNodeInfos, nodeIpInfo)
	}

	serviceNodes = &tmpNodeInfos
}
