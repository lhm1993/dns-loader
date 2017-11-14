package web

import (
	"log"

	"github.com/nu7hatch/gouuid"
	"github.com/zhangmingkai4315/dns-loader/dnsloader"
)

// RPCService define the rpc service object
type RPCService struct {
	ID     *uuid.UUID
	Status Event
	config dnsloader.Configuration
}

// NewRPCService create new service
func NewRPCService() *RPCService {
	id, _ := uuid.NewV4()
	return &RPCService{
		ID:     id,
		Status: Ready,
	}
}

// Ping will response to master ping protocal
func (rpcService *RPCService) Ping(call *RPCCall, result *int) error {
	log.Println(rpcService, call, result)
	// result.ID = *rpcService.ID
	// result.Event = call.Event
	*result = 10
	return nil
}
