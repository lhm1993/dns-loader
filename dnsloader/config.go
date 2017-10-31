package dnsloader

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/ini.v1"
)

// Configuration define all config for this app
type Configuration struct {
	LoaderType         string `json:"loader_type"`
	ControlMaster      string `json:"control_master"`
	Duration           int    `json:"duration"`
	QPS                int    `json:"qps"`
	Server             string `json:"server"`
	Port               int    `json:"port"`
	Domain             string `json:"domain"`
	DomainRandomLength int    `json:"domain_random_length"`
	QueryTypeFixed     bool   `json:"query_type_fixed"`
	QueryType          string `json:"query_type"`
	Debug              bool   `json:"debug"`
	HTTPServer         string `json:"web"`
	RPCPort            int    `json:"rpc_port"`
}

// LoadConfigurationFromIniFile func read a .ini file from file system
// and return the config object
func (config *Configuration) LoadConfigurationFromIniFile(filename string) (err error) {
	if strings.HasSuffix(filename, ".ini") == false {
		return errors.New("Configuration file must be .ini file type")
	}
	cfg, err := ini.Load(filename)
	if err != nil {
		return fmt.Errorf("Configuration file load error:%s", err.Error())
	}
	secApp, err := cfg.GetSection("App")
	if err != nil {
		return fmt.Errorf("Configuration file load section [App] error:%s", err.Error())
	}
	secQuery, err := cfg.GetSection("Query")
	if err != nil {
		return fmt.Errorf("Configuration file load section [Query] error:%s", err.Error())
	}
	// load app attribute
	if secApp.HasKey("type") {
		config.LoaderType = secApp.Key("type").String()
	}
	if secApp.HasKey("control_master") {
		config.ControlMaster = secApp.Key("control_master").String()
	}
	if secApp.HasKey("rpc_port") {
		config.RPCPort = secQuery.Key("rpc_port").MustInt()
	}
	if secApp.HasKey("http_server") {
		config.HTTPServer = secApp.Key("http_server").String()
	}
	if secApp.HasKey("debug") {
		config.Debug = secApp.Key("debug").MustBool()
	}
	// Load traffic attribute
	if secQuery.HasKey("duration") {
		config.Duration = secQuery.Key("duration").MustInt()
	}
	if secQuery.HasKey("qps") {
		config.QPS = secQuery.Key("qps").MustInt()
	}
	if secQuery.HasKey("server") {
		config.Server = secQuery.Key("server").String()
	}
	if secQuery.HasKey("port") {
		config.Port = secQuery.Key("port").MustInt()
	}

	if secQuery.HasKey("domain") {
		config.Domain = secQuery.Key("domain").String()
	}
	if secQuery.HasKey("randomlen") {
		config.DomainRandomLength = secQuery.Key("server").MustInt()
	}
	if secQuery.HasKey("enable_random_type") {
		config.QueryTypeFixed = secQuery.Key("enable_random_type").MustBool()
	}
	if secQuery.HasKey("query_type") {
		config.QueryType = secQuery.Key("query_type").String()
	}

	return
}
