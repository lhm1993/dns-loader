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
	Addr   string
	Port   int
	packet *dns.DNSPacket
	Config *Config
	Conn   net.Conn
}

// NewDefaultConfig func return the default configration object
func NewDefaultConfig() (config *Config) {
	return &Config{
		// 使用本地ip地址
		LocalIP: true,
		// 是否源地址固定
		// SourceIPFixed:false,
		// 是否固定域名
		DomainFixed: false,
		// 固定部分的域名
		Domain: "test1",
		// 随机域名长度
		DomainRandomLength: 1,
		// 是否查询类型固定
		QueryTypeFixed: true,
		// 固定的查询类型
		QueryType: "A",
		// 解析请求的数量
	}
}

// NewDNSClientWithDefaultConfig 使用默认的发包参数进行初始化
func NewDNSClientWithDefaultConfig(addr string, port int) (*DNSClient, error) {
	config := NewDefaultConfig()
	dnsclient, err := NewDNSClientWithConfig(addr, port, config)
	if err != nil {
		return nil, err
	}
	return dnsclient, nil
}

// NewDNSClientWithConfig 使用传递的配置发包参数进行初始化
func NewDNSClientWithConfig(addr string, port int, config *Config) (dnsclient *DNSClient, err error) {
	if port > 65535 || port < 0 {
		port = 53
	}
	dnsclient = &DNSClient{Addr: addr, Port: port, Config: config}
	// udpAddr, err := net.ResolveUDPAddr("udp", addr+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("udp", addr+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	dnsclient.Conn = conn
	log.Println("New DNS client success, start send packet...")
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
	log.Println("Init DNS client configuration success")
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
		log.Printf("Send DNS Query Failed:%s", err)
		return
	}
}
