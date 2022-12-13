package cores

import "apioak-admin/routers"

func InitRouter(conf *ConfigGlobal) error {
	routers.RouterRegister(conf.Runtime.Gin)
	return nil
}
