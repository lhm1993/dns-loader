package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/zhangmingkai4315/dns-loader/dnsloader"
	"github.com/zhangmingkai4315/dns-loader/web"
)

var (
	loaderType string
	timeout    int
	master     string
	duration   int
	max        int
	qps        int
	domain     string
	server     string
	port       int
	randomlen  int
	randomtype bool
	configFile string
	queryType  string
	debug      bool
)

func init() {
	flag.StringVar(&loaderType, "t", "once", "loader type[master|worker|once(default)]")
	flag.StringVar(&configFile, "f", "", "config file")
	flag.StringVar(&master, "master", "", "the master host ip and port 192.168.1.1:8000")
	flag.IntVar(&duration, "duration", 60, "the duration time (s)")
	flag.IntVar(&timeout, "timeout", 1, "the timeout for one query[not implement]")
	flag.IntVar(&qps, "q", 10, "dns query per second(0=unlimit speed)")
	flag.StringVar(&server, "s", "172.17.0.2", "dns server and listen port")
	flag.IntVar(&port, "p", 53, "dns query server port")
	flag.StringVar(&domain, "domain", "google.com", "base domain,for example :google.com")
	flag.IntVar(&randomlen, "len", 5, "random length of subdomain, for example 5 means *****.google.com")
	flag.BoolVar(&randomtype, "randomtype", false, "random dns type to send")
	flag.StringVar(&queryType, "querytype", "A", "dns query type")
	flag.BoolVar(&debug, "d", false, "enable debug mode")
}

func main() {
	flag.Parse()
	//if config file exist, load config file
	var config *dnsloader.Configuration
	config = &dnsloader.Configuration{}
	if configFile != "" {
		log.Printf("load configuration from file:%s\n", configFile)
		err := config.LoadConfigurationFromIniFile(configFile)
		if err != nil {
			log.Panicf("read configuration file error:%s", err.Error())
		}
	} else {
		log.Println("load configuration from command line")
		config.ControlMaster = master
		config.LoaderType = loaderType
		config.Domain = domain
		config.DomainRandomLength = randomlen
		config.QPS = qps
		config.Duration = duration
		config.QueryTypeFixed = randomtype
		config.Server = server
		config.Port = port
		config.QueryType = queryType
		config.Debug = debug
	}
	// enable performance debug for app
	if config.Debug == true {
		go func() {
			log.Println("Start performace monitor on port 8080")
			err := http.ListenAndServe("localhost:8080", http.DefaultServeMux)
			if err != nil {
				log.Println("Start performance monitor fail")
			}
		}()
	}
	// Start docker container first
	if config.LoaderType == "once" {
		dnsloader.GenTrafficFromConfig(config)
	} else if config.LoaderType == "master" {
		// start web component
		if configFile == "" {
			log.Fatalln("Please using -f to load config file first")
		}
		log.Printf("Start Web for control panel default web address:%s\n", config.HTTPServer)
		web.NewServer(config)
	} else if config.LoaderType == "agent" {
		// start rpc regist to server
		
		// start rpc waiting status
	}
}
