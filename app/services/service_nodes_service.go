package services

import (
	"apioak-admin/app/models"
	"apioak-admin/app/utils"
	"apioak-admin/app/validators"
	"strconv"
)

func GetToOperateNodes(serviceId string, serviceNodes *[]validators.ServiceNodeAddUpdate) ([]models.ServiceNodes, []models.ServiceNodes, []string) {
	serviceNodesModel := models.ServiceNodes{}
	serviceExistNodes := serviceNodesModel.NodeInfosByServiceIds([]string{serviceId})

	updateNodesMap := make(map[string]validators.ServiceNodeAddUpdate)
	for _, updateNode := range *serviceNodes {
		nodePort := strconv.Itoa(updateNode.NodePort)
		updateNodeValue := updateNode.NodeIp + "-" + nodePort
		updateNodesMap[updateNodeValue] = updateNode
	}

	existNodesMap := make(map[string]string)
	for _, serviceExistNode := range serviceExistNodes {
		nodePort := strconv.Itoa(serviceExistNode.NodePort)
		existNodeValue := serviceExistNode.NodeIP + "-" + nodePort
		existNodesMap[existNodeValue] = existNodeValue
	}

	addNodes := make([]models.ServiceNodes, 0)
	for _, addNode := range *serviceNodes {
		nodePort := strconv.Itoa(addNode.NodePort)
		updateNodeValue := addNode.NodeIp + "-" + nodePort
		_, exist := existNodesMap[updateNodeValue]
		if exist {
			continue
		}

		ipType, err := utils.DiscernIP(addNode.NodeIp)
		if err != nil {
			continue
		}
		ipTypeMap := models.IPTypeMap()
		nodeInfo := models.ServiceNodes{
			ServiceID:  serviceId,
			NodeIP:     addNode.NodeIp,
			IPType:     ipTypeMap[ipType],
			NodePort:   addNode.NodePort,
			NodeWeight: addNode.NodeWeight,
		}
		addNodes = append(addNodes, nodeInfo)
	}

	updateNodes := make([]models.ServiceNodes, 0)
	deleteNodeIds := make([]string, 0)
	for _, serviceExistNode := range serviceExistNodes {
		nodePort := strconv.Itoa(serviceExistNode.NodePort)
		existNodeValue := serviceExistNode.NodeIP + "-" + nodePort
		updateNode, exist := updateNodesMap[existNodeValue]
		if exist {
			if updateNode.NodeWeight == serviceExistNode.NodeWeight {
				continue
			}
			nodeInfo := models.ServiceNodes{
				ID:         serviceExistNode.ID,
				ServiceID:  serviceId,
				NodeIP:     updateNode.NodeIp,
				IPType:     serviceExistNode.IPType,
				NodePort:   updateNode.NodePort,
				NodeWeight: updateNode.NodeWeight,
			}
			updateNodes = append(updateNodes, nodeInfo)
		} else {
			deleteNodeIds = append(deleteNodeIds, serviceExistNode.ID)
		}
	}

	return addNodes, updateNodes, deleteNodeIds
}
