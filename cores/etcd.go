package cores

import (
	"apioak-admin/app/packages"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

func InitEtcd(conf *ConfigGlobal) error {
	if len(conf.Etcd.HostPort) <= 0 {
		return errors.New("etcd configuration cannot be empty")
	}

	hostPorts := strings.Split(conf.Etcd.HostPort, ",")
	if len(hostPorts) <= 0 {
		return fmt.Errorf("etcd configuration format error: `%s`", conf.Etcd.HostPort)
	}

	for _, hostPortStr := range hostPorts {
		hostPortStruct := strings.Split(hostPortStr, ":")
		if len(hostPortStruct) != 2 {
			return fmt.Errorf("etcd configuration format error: `%s`", conf.Etcd.HostPort)
		}
	}

	client, clientErr := clientv3.New(clientv3.Config{
		Endpoints:   hostPorts,
		DialTimeout: time.Second * 2,
	})
	if clientErr != nil {
		return fmt.Errorf("etcd client initialization error: `%s`", clientErr.Error())
	}

	packages.SetEtcdClient(client)

	return nil
}
