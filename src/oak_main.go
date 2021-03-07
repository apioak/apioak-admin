package src

import (
	"log"
)

func Main() {
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	//cli.HelpFlag = cli.BoolFlag{
	//	Name: "help, h",
	//	Usage: "show help and exit",
	//}
	//
	//cli.VersionFlag = cli.BoolFlag{
	//	Name: "version, v",
	//	Usage: "show version and exit",
	//}

	//if err := initConfig(); err != nil {
	//	panic(err)
	//}
	//
	//if config.CLI.Version {
	//	fmt.Printf("APIOAK: Version %s\n", AppVersion)
	//	os.Exit(1)
	//}
	//
	//if len(config.Etcd.Nodes) == 0 {
	//	fmt.Print("error: config etcd nodes is empty")
	//	os.Exit(1)
	//}
	//app := gin.Default()
	//app.Run(":8080")
}
