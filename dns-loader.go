package main

import (
	"flag"
	"github.com/zhangmingkai4315/dns-loader/dnsloader"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
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
	flag.StringVar(&loaderType, "type", "master", "loader type[master|worker]")
	flag.StringVar(&master, "master", "", "the master host")
	flag.IntVar(&duration, "duration", 60, "the duration time (s)")
	flag.IntVar(&timeout, "timeout", 1, "the timeout for one query[not implement]")
	flag.IntVar(&qps, "qps", 10, "dns query per second(0=unlimit speed)")
	flag.StringVar(&server, "server", "172.17.0.2", "dns server and listen port")
	flag.IntVar(&port, "port", 53, "dns query server port")
	flag.StringVar(&domain, "domain", "google.com", "base domain,for example :google.com")
	flag.IntVar(&randomlen, "len", 5, "random length of subdomain, for example 5 means *****.google.com")
	flag.BoolVar(&randomtype, "randomtype", false, "random dns type to send")
	flag.StringVar(&queryType, "querytype", "A", "dns query type")
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
}

func main() {
	flag.Parse()
	//if config file exist, load config file
	var config *dnsloader.Configuration
	if configFile != "" {
		log.Printf("load configuration from file:%s\n", configFile)
		err := config.LoadConfigurationFromIniFile(configFile)
		if err != nil {
			log.Panicf("read configuration file error:%s", err.Error())
		}
	} else {
		config = &dnsloader.Configuration{}
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

	dnsclient, err := dnsloader.NewDNSClientWithConfig(config)
	if err != nil {
		log.Panicf("%s", err.Error())
	}
	log.Println("config the dns loader success")
	log.Printf("current configuration for dns loader is server:%s|port:%d\n",
		dnsclient.Config.Server, dnsclient.Config.Port)
	// log.Printf("%+v", config)
	param := dnsloader.GeneratorParam{
		Caller:   dnsclient,
		Timeout:  1000 * time.Millisecond,
		QPS:      uint32(config.QPS),
		Duration: time.Second * time.Duration(config.Duration),
	}
	log.Printf("initialize load %+v", param)
	gen, err := dnsloader.NewDNSLoaderGenerator(param)
	if err != nil {
		log.Panicf("load generator initialization fail :%s", err)
	}
	log.Println("start load generator")
	gen.Start()
}
