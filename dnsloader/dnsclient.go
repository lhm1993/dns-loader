package dnsloader

import (
	// "bytes"
	"github.com/zhangmingkai4315/go-dns-shooter/dns"
	"log"
	"net"
	"strconv"
)

// DNSClient 定义dnsloader发包的配置数据
type DNSClient struct {
	packet *dns.DNSPacket
	Config *Configuration
	Conn   net.Conn
}

// NewDNSClientWithConfig 使用传递的配置发包参数进行初始化
func NewDNSClientWithConfig(config *Configuration) (dnsclient *DNSClient, err error) {
	dnsclient = &DNSClient{Config: config}
	conn, err := net.Dial("udp", config.Server+":"+strconv.Itoa(config.Port))
	if err != nil {
		return nil, err
	}
	dnsclient.Conn = conn
	log.Println("new dns loader client success, start send packet...")
	dnsclient.InitPacket()
	return dnsclient, nil
}

// InitPacket init a packet for dns query data
func (client *DNSClient) InitPacket() {
	client.packet = new(dns.DNSPacket)
	// client.packet.InitialPacket(domain, randomlen, dns.TypeA)
	if client.Config.QueryTypeFixed == true {
		client.packet.InitialPacket(client.Config.Domain,
			client.Config.DomainRandomLength,
			GetDNSTypeCodeFromString(client.Config.QueryType))
	} else {
		client.packet.InitialPacket(client.Config.Domain,
			client.Config.DomainRandomLength,
			dns.TypeA)
		client.packet.RandomType = true
	}
	log.Println("init dns loader client configuration success")
}

// BuildReq build new dns request for use later
func (client *DNSClient) BuildReq() []byte {
	randomDomain := dns.GenRandomDomain(client.Config.DomainRandomLength, client.Config.Domain)
	if _, err := client.packet.UpdateSubDomainToBytes(randomDomain); err != nil {
		log.Printf("%v\n", err)
	}
	return client.packet.RawByte
}

// Call func will be called by schedual each time
func (client *DNSClient) Call(req []byte) {
	_, err := client.Conn.Write(req)
	if err != nil {
		log.Printf("send dns query Failed:%s", err)
		return
	}
}
