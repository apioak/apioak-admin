package cores

import (
	"apioak-admin/app/services"
	"time"
)

func InitGoroutineFunc() {
	// go ClusterNodeWatch()

	go dynamicValidationPluginData()

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

func dynamicValidationPluginData() {

	timer := time.NewTicker(10 * time.Second)
	defer timer.Stop()

	for {
		services.PluginBasicInfoMaintain()

		<-timer.C

	}
}
