package packages

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

var (
	dbConnection *gorm.DB
	etcdConnection *clientv3.Client
)

func SetDb(db *gorm.DB) {
	dbConnection = db
}

func GetDb() *gorm.DB {
	return dbConnection
}

func SetEtcdClient(etcdClient *clientv3.Client) {
	etcdConnection = etcdClient
}

func GetEtcdClient() *clientv3.Client {
	return etcdConnection
}
