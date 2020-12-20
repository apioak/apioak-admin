package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

func startServer() {
	u := &Upstream{"127.0.0.1", "10222"}
	err := http.ListenAndServe(":10111", u)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	fmt.Println("Hello, APIOAK")
	startServer()
}
