package packages

var pluginKeys []string

func SetPluginKeys(pluginKeyList []string) (err error) {

	pluginKeys = pluginKeyList

	return
}

func GetPluginKeys() []string {
	return pluginKeys
}

