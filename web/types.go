package web

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
	"github.com/zhangmingkai4315/dns-loader/dnsloader"
)

// Event define the event send to all node
type Event uint8

const (
	// Start and send the new config
	Start Event = iota
	// Check the status
	Check
	// Running status with some message
	Running
	// Kill the load
	Kill
	// Error Status
	Error
	// Stop the load in normal way
	Stop
	// Ping will do health check the status of node
	Ping
)

// IPWithPort define the posted node info
type IPWithPort struct {
	IPAddress string `json:"ipaddress"`
	Port      int    `json:"port"`
}

func (ipp *IPWithPort) toString(defaultPort int) string {
	if ipp.Port == 0 {
		ipp.Port = defaultPort
	}
	return fmt.Sprintf("%s:%d", ipp.IPAddress, ipp.Port)
}

// RPCCall define the message send to node
type RPCCall struct {
	ID      uuid.UUID
	Event   Event
	Config  dnsloader.Configuration
	Message string
}

// RPCResult define the result send from node
type RPCResult struct {
	ID      uuid.UUID
	Event   Event
	Message string
}
