package cores

import (
	"apioak-admin/app/services"
	"time"
)

func InitGoroutineFunc() {
	go dynamicValidationPluginData()
}

func dynamicValidationPluginData() {

	timer := time.NewTicker(10 * time.Second)
	defer timer.Stop()

	for {
		services.PluginBasicInfoMaintain()

		<-timer.C
	}
}
