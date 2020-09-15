package server

import (
	"log"
	"net"

	"github.com/spf13/viper"
)

var dnsServer *DNSServer

type DNSServer struct {
	DomainMap map[string]string
	BlockedIP map[string]interface{}
	socket    *net.UDPConn
}

func GetDNSServer() *DNSServer {
	if dnsServer != nil {
		return dnsServer
	}
	var err error
	dnsServer = &DNSServer{
		DomainMap: make(map[string]string),
		BlockedIP: make(map[string]interface{}),
	}
	dnsServer.DomainMap = viper.GetStringMapString("domain_map")
	blockedIPs := viper.GetStringSlice("blocked_ip")
	for index := range blockedIPs {
		dnsServer.BlockedIP[blockedIPs[index]] = struct{}{}
	}

	dnsServer.socket, err = net.ListenUDP(UDP_NETWORK, &net.UDPAddr{
		IP:   net.ParseIP(viper.GetString("dns_relay.client_ip")),
		Port: DNS_PORT,
	})
	if err != nil {
		log.Printf("配置错误：%v", err)
		return dnsServer
	}

	log.Printf("本地共%v条数据\n", len(dnsServer.DomainMap))
	for k, v := range dnsServer.DomainMap {
		log.Printf("域名：%v，对应IP：%v", k, v)
	}
	log.Printf("本地共%v条屏蔽ip\n", len(dnsServer.BlockedIP))
	for k := range dnsServer.BlockedIP {
		log.Printf("屏蔽ip：%v", k)
	}
	return dnsServer
}

func (dnsServer *DNSServer) Serve() {
	for {
		data := make([]byte, 1024)
		read, remoteAddr, err := dnsServer.socket.ReadFromUDP(data)
		if err != nil {
			log.Println("接收数据错误", err)
			continue
		}
		data = data[:read]
		log.Println("接收数据成功", read, remoteAddr, data)
		parserServer, err := GetParserServer(data, remoteAddr)
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
