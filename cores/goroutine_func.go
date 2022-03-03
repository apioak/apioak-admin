package cores

func InitGoroutineFunc() {
	go ClusterNodeWatch()
}

func ClusterNodeWatch() {
	//etcdClient := packages.GetEtcdClient()
	//for true {
	//	rch := etcdClient.Watch(context.TODO(), utils.EtcdKeyWatchClusterNode)
	//	for wresp := range rch {
	//		for _, ev := range wresp.Events {
	//
	//			if ev.Type.String() == "PUT" {
	//				services.ClusterNodeWatchAdd(string(ev.Kv.Value))
	//			}
	//		}
	//	}
	//}
}
