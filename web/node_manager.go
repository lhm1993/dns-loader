package web

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"

	"github.com/zhangmingkai4315/dns-loader/dnsloader"

	"github.com/nu7hatch/gouuid"
)

// NodeStatus define the one of node status
type NodeStatus struct {
	Status    Event
	TimeStamp time.Time
	Message   string
}

// NewNodeStatus create the status from message
func NewNodeStatus(status Event, message string) *NodeStatus {
	return &NodeStatus{
		Status:    status,
		Message:   message,
		TimeStamp: time.Now(),
	}
}

// NodeManager define the rpc node list
// when new config generated the manager will call the nodes one by one
type NodeManager struct {
	mutex      *sync.RWMutex
	IPList     []string
	IPStatus   map[string][]NodeStatus
	TaskStatus []string
	config     *dnsloader.Configuration
}

// NewNodeManager will create a new node manager
func NewNodeManager(c *dnsloader.Configuration) *NodeManager {
	ipstatus := make(map[string][]NodeStatus)
	return &NodeManager{
		mutex:      &sync.RWMutex{},
		IPList:     []string{},
		TaskStatus: []string{},
		IPStatus:   ipstatus,
		config:     c,
	}
}

// AddNode append a new ip to this node list
func (manager *NodeManager) AddNode(ip string, port int) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	if port == 0 {
		port = manager.config.RPCPort
	}
	ip = fmt.Sprintf("%s:%d", ip, port)
	log.Printf("send rpc ping request to node:%s\n", ip)
	err := manager.callPing(ip, Ping)
	if err != nil {
		return err
	}
	manager.IPList = append(manager.IPList, ip)
	return nil
}

// AddStatus append a new ip to this node list
func (manager *NodeManager) AddStatus(ip string, status Event, message string) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	statuslist, ok := manager.IPStatus[ip]
	if ok != true {
		return errors.New("ip not exist")
	}
	nodeStatus := NewNodeStatus(status, message)
	statuslist = append(statuslist, *nodeStatus)
	return nil
}

// Remove will remove the ip from current list
func (manager *NodeManager) Remove(deleteip string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	for i, ip := range manager.IPList {
		if ip == deleteip {
			manager.IPList = append(manager.IPList[:i], manager.IPList[i+1:]...)
			break
		}
	}
	delete(manager.IPStatus, deleteip)
}

// Call function will send data to all node
func (manager *NodeManager) Call(event Event, data interface{}) error {

	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, ip := range manager.IPList {
		go func(ip string, event Event, data interface{}) {
			switch event {
			case Start:
				manager.callStart(ip, event, data)
			case Check:
				manager.callCheckStatus(ip, event, data)
			case Kill:
				manager.callKill(ip, event, data)
			case Ping:
				manager.callPing(ip, event)
			}
		}(ip, event, data)
	}
	return nil
}

// callStart function will send data to all node with start evnet and config
func (manager *NodeManager) callStart(ip string, event Event, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic from rpc call[%s]\n", ip)
		}
	}()
	client, err := rpc.DialHTTP("tcp", ip)
	defer client.Close()
	if err != nil {
		log.Printf("call remote node rpc failed:[%s] %s\n", ip, err.Error())
		return
	}
	config, ok := data.(dnsloader.Configuration)
	if ok != true {
		log.Println("config data fail to send")
		return
	}
	var result *RPCResult
	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("uuid generate failed :%s \n", err.Error())
		return
	}
	args := &RPCCall{
		Event:  event,
		Config: config,
		ID:     *id,
	}
	err = client.Call("Manager.GenDNSTraffic", args, &result)
	if err != nil {
		log.Printf("send configuration to node faile:[%s] %s", ip, err.Error())
		return
	}
	return nil
}

// callCheckStatus function will check all the node with uuid
func (manager *NodeManager) callCheckStatus(ip string, event Event, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic from rpc call[%s]\n", ip)
		}
	}()
	client, err := rpc.DialHTTP("tcp", ip)
	defer client.Close()
	if err != nil {
		log.Printf("call remote node rpc failed:[%s] %s\n", ip, err.Error())
		return
	}
	id, ok := data.(uuid.UUID)
	if ok != true {
		log.Println("config data fail to send")
		return
	}
	var result *RPCResult
	args := &RPCCall{
		Event: event,
		ID:    id,
	}
	err = client.Call("Manager.CheckStatus", args, &result)
	if err != nil {
		log.Printf("send check status to node faile:[%s] %s", ip, err.Error())
		return
	}
	return nil
}

// callCheckStatus function will check all the node with uuid
func (manager *NodeManager) callKill(ip string, event Event, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic from rpc call[%s]\n", ip)
		}
	}()
	client, err := rpc.DialHTTP("tcp", ip)
	defer client.Close()
	if err != nil {
		log.Printf("call remote node rpc failed:[%s] %s\n", ip, err.Error())
		return
	}
	id, ok := data.(uuid.UUID)
	if ok != true {
		log.Println("config data fail to send")
		return
	}
	var result *RPCResult
	args := &RPCCall{
		Event: event,
		ID:    id,
	}
	err = client.Call("Manager.KillProcess", args, &result)
	if err != nil {
		log.Printf("send check status to node fail:[%s] %s", ip, err.Error())
		return
	}
	return nil
}

// callCheckStatus function will check all the node with uuid
func (manager *NodeManager) callPing(ip string, event Event) (err error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Printf("panic from rpc call[%s]\n", r)
	// 	}
	// }()
	client, err := jsonrpc.Dial("tcp", ip)
	defer client.Close()
	if err != nil {
		log.Printf("call remote node rpc fail:[%s] %s\n", ip, err.Error())
		return
	}
	var result int
	args := RPCCall{
		Event: event,
	}
	err = client.Call("RPCService.Ping", args, &result)
	if err != nil {
		log.Printf("send ping check node fail:[%s] %s", ip, err.Error())
		return
	}
	log.Printf("send ping check success:[%s] %s", ip, err.Error())
	return nil
}
