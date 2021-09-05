package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strings"
)

func createRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func IdGenerate(idType string) (string, error) {
	randomId := createRandomString(IdLength)

	var id string
	switch strings.ToLower(idType) {
	case IdTypeUser:
		id = IdTypeUser + "-" + randomId
	case IdTypeService:
		id = IdTypeService + "-" + randomId
	case IdTypeServiceDomain:
		id = IdTypeServiceDomain + "-" + randomId
	case IdTypeServiceNode:
		id = IdTypeServiceNode + "-" + randomId
	case IdTypeRoute:
		id = IdTypeRoute + "-" + randomId
	case IdTypeRoutePlugin:
		id = IdTypeRoutePlugin + "-" + randomId
	case IdTypeCertificate:
		id = IdTypeCertificate + "-" + randomId
	case IdTypeClusterNode:
		id = IdTypeClusterNode + "-" + randomId
	default:
		return "", fmt.Errorf("id type error")
	}

	return id, nil
}

func DiscernIP(s string) (string, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return "", fmt.Errorf("(%s) is illegal ip", s)
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return IPV4, nil
		case ':':
			return IPV6, nil
		}
	}
	return "", nil
}

type LoadBalance struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (l *LoadBalance) LoadBalanceList() []LoadBalance {
	loadBalance := LoadBalance{}
	loadBalanceList := make([]LoadBalance, 0)

	loadBalance.Id = LoadBalanceRoundRobin
	loadBalance.Name = LoadBalanceNameRoundRobin
	loadBalanceList = append(loadBalanceList, loadBalance)

	loadBalance.Id = LoadBalanceIPHash
	loadBalance.Name = LoadBalanceNameIPHash
	loadBalanceList = append(loadBalanceList, loadBalance)

	return loadBalanceList
}
