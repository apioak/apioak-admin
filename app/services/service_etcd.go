package services

import (
	"apioak-admin/app/packages"
	"apioak-admin/app/utils"
	"context"
	"errors"
	"time"
)

func EtcdSync(key string, value string) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*utils.EtcdTimeOut)
	etcdClient := packages.GetEtcdClient()
	_, err := etcdClient.Put(ctx, key, value)
	cancel()

	if err != nil {
		return false, err
	}

	return true, nil
}

func EtcdGet(key string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*utils.EtcdTimeOut)
	etcdClient := packages.GetEtcdClient()
	etcdValue, etcdValueErr := etcdClient.Get(ctx, key)
	cancel()

	if etcdValueErr != nil {
		return "", etcdValueErr
	}

	if len(etcdValue.Kvs) == 0 {
		return "", errors.New("data is empty")
	}

	return string(etcdValue.Kvs[0].Value), nil
}
