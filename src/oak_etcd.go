package src

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

type OakEtcd struct {
	config ConfigEtcd
	client *clientv3.Client
}

func (oe *OakEtcd) New() (*OakEtcd, error) {
	if len(oe.config.Nodes) == 0 {
		return nil, errors.New("etcd config nodes is empty")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: oe.config.Nodes,
	})

	if err != nil {
		return nil, err
	}

	defer cli.Close()

	oe.client = cli

	return oe, nil
}

func (oe *OakEtcd) WatchPrefix(prefix string) {
	for true {
		rch := oe.client.Watch(context.Background(), prefix)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}

func initEtcd(config ConfigEtcd) (*OakEtcd, error) {
	oakEtcd := OakEtcd{
		config: config,
	}

	return oakEtcd.New()
}
