package dnsloader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nu7hatch/gouuid"

	"github.com/asaskevich/govalidator"
	"gopkg.in/ini.v1"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Configuration define all config for this app
type Configuration struct {
	ID                 string `json:"id" valid:"uuid"`
	LoaderType         string `json:"loader_type" valid:"in(once|master|agent),optional"`
	ControlMaster      string `json:"control_master" valid:"ip,optional"`
	Duration           int    `json:"duration" valid:"-"`
	QPS                int    `json:"qps" valid:"-"`
	Server             string `json:"server" valid:"ip"`
	Port               int    `json:"port" valid:"-"`
	Domain             string `json:"domain" valid:"dns"`
	DomainRandomLength int    `json:"domain_random_length" valid:"-"`
	QueryTypeFixed     bool   `json:"query_type_fixed" valid:"-"`
	QueryType          string `json:"query_type" valid:"-"`
	Debug              bool   `json:"debug" valid:"-"`
	HTTPServer         string `json:"web" valid:"ip,optional"`
	AgentPort          int    `json:"agent_port" valid:"-"`

	User       string `json:"-" valid:"-"`
	Password   string `json:"-" valid:"-"`
	AppSecrect string `json:"-" valid:"-"`
}

// Valid will check all setting
func (config *Configuration) Valid() error {
	if config.ID == "" {
		id, _ := uuid.NewV4()
		config.ID = (*id).String()
	}
	_, err := govalidator.ValidateStruct(config)
	if err != nil {
		return err
	}
	if config.QPS < 0 ||
		config.Duration < 0 ||
		config.DomainRandomLength < 0 ||
		config.AgentPort < 0 ||
		config.Port < 0 {
		return errors.New("number can't set to nagetive")
	}
	return nil
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

	// load app attribute
	if secApp.HasKey("user") {
		config.User = secApp.Key("user").String()
	}
	if secApp.HasKey("password") {
		config.Password = secApp.Key("password").String()
	}
	if secApp.HasKey("app_secrect") {
		config.AppSecrect = secApp.Key("app_secrect").String()
	}
	if secApp.HasKey("control_master") {
		config.ControlMaster = secApp.Key("control_master").String()
	}
	if secApp.HasKey("agent_port") {
		config.AgentPort = secApp.Key("agent_port").MustInt()
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
