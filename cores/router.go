package cores

import "apioak-admin/routes"

func InitRoute(conf *ConfigGlobal) error {
	routes.AdminRegister(conf.Runtime.Gin)
	return nil
}
