package cores

import "apioak-admin/routes"

func InitRoute(conf *ConfigGlobal) error {
	routes.RouteRegister(conf.Runtime.Gin)
	return nil
}
