package server

import (
	"log"
	"net"
	"sync"

	"github.com/spf13/viper"
)

var dnsServer *DNSServer
var once *sync.Once

type DNSServer struct {
	DomainMap map[string]string
	socket    *net.UDPConn
}

func GetDNSServer() *DNSServer {
	once.Do(func() {
		var err error
		dnsServer = &DNSServer{}
		dnsServer.DomainMap = viper.GetStringMapString("domain_map")
		clientIP := viper.GetIntSlice("dns_relay.client_ip")
		dnsServer.socket, err = net.ListenUDP(UDP_NETWORK, &net.UDPAddr{
			IP:   net.IPv4(byte(clientIP[0]), byte(clientIP[1]), byte(clientIP[2]), byte(clientIP[3])),
			Port: viper.GetInt("dns_relay.client_port"),
		})
		if err != nil {
			log.Panicf("配置错误：%v", err)
		}

		log.Printf("本地共%v条数据\n", len(dnsServer.DomainMap))
		log.Println(dnsServer.DomainMap)

	})
	return dnsServer
}

func (dnsServer *DNSServer) Serve() {
	data := make([]byte, 1024)
	for {
		read, remoteAddr, err := dnsServer.socket.ReadFromUDP(data)
		if err != nil {
			log.Println("接收数据错误", err)
			continue
		}
		parserServer, err := GetParserServer(data, remoteAddr)
		log.Println("接收数据成功", read, remoteAddr, data)
		if err != nil {
			log.Println(err)
			continue
		}
		go parserServer.parse() //nolint: errcheck
	}
}

func (dnsServer *DNSServer) Close() {
	if dnsServer.socket != nil {
		dnsServer.socket.Close()
	}
}
