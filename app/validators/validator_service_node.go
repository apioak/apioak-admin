package validators

import (
	"apioak-admin/app/packages"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ServiceNodeAddUpdate struct {
	NodeIp     string `json:"node_ip" zh:"上游节点IP" en:"Node IP" binding:"required,ip"`
	NodePort   int    `json:"node_port" zh:"端口" en:"Node port" binding:"omitempty,min=1,max=65535"`
	NodeWeight int    `json:"node_weight" zh:"权重" en:"Node weight" binding:"omitempty,min=1,max=100"`
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
