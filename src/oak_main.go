package src

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

type Upstream struct {
	host string
	port string
}

func (upstream *Upstream) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	remote, err := url.Parse("http://" + upstream.host + ":" + upstream.port)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(writer, request)
}

func startServer(application *ConfigApp) {
	u := &Upstream{"127.0.0.1", "10222"}
	err := http.ListenAndServe(application.Host+":"+strconv.Itoa(application.Port), u)
	if err != nil {
		log.Fatalln(err)
	}
}

func Main() {
	flag.Parse()

	if err := initConfig(); err != nil {
		panic(err)
	}

	if config.CLI.Version {
		fmt.Printf("APIOAK: Version %s\n", AppVersion)
		os.Exit(1)
	}

	if len(config.Etcd.Nodes) == 0 {
		fmt.Print("error: config etcd nodes is empty")
		os.Exit(1)
	}

	startServer(&config.App)
}
