package main

import (
	"apioak-admin/cores"
)

func main() {

	// 全局配置
	var conf cores.ConfigGlobal

	// 全局配置初始化
	if err := cores.InitConfig(&conf); err != nil {
		panic(err)
	}

	// 初始化框架
	if err := cores.InitFramework(&conf); err != nil {
		panic(err)
	}

	// 初始化Logger
	if err := cores.InitLogger(&conf); err != nil {
		panic(err)
	}

	// 初始化数据库
	if err := cores.InitDataBase(&conf); err != nil {
		panic(err)
	}

	//初始化ETCD
	if err := cores.InitEtcd(&conf); err != nil {
		panic(err)
	}

	//初始化Token
	if err := cores.InitToken(&conf); err != nil {
		panic(err)
	}

	// 初始化参数验证器
	if err := cores.InitValidator(&conf); err != nil {
		panic(err)
	}

	// 初始化路由
	if err := cores.InitRoute(&conf); err != nil {
		panic(err)
	}

	// 协程处理额外事件
	// cores.InitGoroutineFunc()

	// 服务启动
	if err := cores.RunFramework(&conf); err != nil {
		panic(err)
	}
}
