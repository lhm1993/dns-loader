package main

import (
	"github.com/zhangmingkai4315/dns-loader/dnsloader"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		log.Println("Start performace monitor on port 8080")
		err := http.ListenAndServe("localhost:8080", http.DefaultServeMux)
		if err != nil {
			log.Println("Start performance monitor fail")
		}
	}()
	// Start docker container first
	server := "172.17.0.2"
	port := 53
	dnsclient, err := dnsloader.NewDNSClientWithDefaultConfig(server, port)
	if err != nil {
		log.Panicf("%s", err.Error())
	}
	log.Println("Config the dns loader done")
	log.Printf("Current configuration for dns loader is server:%s|port:%d|Detail:%+v\n",
		dnsclient.Addr, dnsclient.Port, dnsclient.Config)

	param := dnsloader.GeneratorParam{
		Caller:   dnsclient,
		Timeout:  1000 * time.Millisecond,
		QPS:      uint32(100000),
		Duration: 20 * time.Second,
		Workers:  2,
	}
	log.Printf("Initialize load %+v", param)
	gen, err := dnsloader.NewDNSLoaderGenerator(param)
	if err != nil {
		log.Panicf("Load generator initialization fail :%s", err)
	}
	log.Println("Start load generator")

	gen.Start()

}
