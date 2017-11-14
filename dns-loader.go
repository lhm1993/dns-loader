package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/zhangmingkai4315/dns-loader/dnsloader"
	"github.com/zhangmingkai4315/dns-loader/web"
)

var (
	loaderType = flag.String("t", "", "")
	timeout    = flag.Int("timeout", 5, "")
	master     = flag.String("master", "", "")
	duration   = flag.Int("duration", 60, "")
	max        = flag.Int("max", 0, "")
	qps        = flag.Int("q", 10, "")
	domain     = flag.String("d", "", "")
	server     = flag.String("s", "", "")
	port       = flag.Int("p", 53, "")
	randomlen  = flag.Int("r", 5, "")
	randomtype = flag.Bool("R", false, "")
	configFile = flag.String("c", "", "")
	queryType  = flag.String("Q", "A", "")
	debug      = flag.Bool("debug", false, "")
)
var usage = `Usage: dns-loader [options...] 
Options:
  -t  loader type, one of "master","worker","once", also can set it in you config file
  -m  master control server ip. Must set this value when type is worker
  -s  dns server. Default "127.0.0.1"
  -p  dns server listen port. Default is 53.
  -d  query domain name
  -D  duration time. Default 60 seconds
  -q  query per second. Default is 10
  -r  random subdomain length. Default is 5
  -R  enable random query type. Default is false
  -c  config file path for app start
  -Q  query type. Default is A
  -debug enable debug mode
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()
	//if config file exist, load config file
	var config *dnsloader.Configuration
	config = &dnsloader.Configuration{}
	// enable performance debug for app
	if *debug == true {
		go func() {
			log.Println("Start performace monitor on port 8080")
			err := http.ListenAndServe("localhost:8080", http.DefaultServeMux)
			if err != nil {
				log.Println("Start performance monitor fail")
			}
		}()
	}

	// Start docker container first
	config.LoaderType = *loaderType
	if *loaderType == "once" {
		config.ControlMaster = *master
		config.Domain = *domain
		config.DomainRandomLength = *randomlen
		config.QPS = *qps
		config.Duration = *duration
		config.QueryTypeFixed = *randomtype
		config.Server = *server
		config.Port = *port
		config.QueryType = *queryType
		config.Debug = *debug
		if err := config.Valid(); err != nil {
			usageAndExit(err.Error())
		}
		dnsloader.GenTrafficFromConfig(config)
		return
	}
	if *configFile != "" {
		log.Printf("load configuration from file:%s\n", *configFile)
		err := config.LoadConfigurationFromIniFile(*configFile)
		if err != nil {
			log.Panicf("read configuration file error:%s", err.Error())
		}
	} else {
		usageAndExit("please using -c to load config file first")
	}
	if *loaderType == "master" {
		log.Printf("start Web for control panel default web address:%s\n", config.HTTPServer)
		web.NewServer(config)
		return
	}

	if *loaderType == "agent" {
		if config.RPCPort == 0 || config.ControlMaster == "" {
			usageAndExit("RPC port and master ip must given")
		}
		log.Printf("start RPC server listen on %d for master:%s connect\n", config.RPCPort, config.ControlMaster)
		rpcService := web.NewRPCService()
		rpc.Register(rpcService)
		tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", config.RPCPort))
		checkError(err)
		listener, err := net.ListenTCP("tcp", tcpAddr)
		checkError(err)
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			jsonrpc.ServeConn(conn)
		}
	}

}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
