package server

import (
	"log"
	"net"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var dnsServer *DNSServer

type DNSServer struct {
	BlockedIP   sync.Map
	RedisClient redis.UniversalClient
	socket      *net.UDPConn
}

func GetDNSServer() *DNSServer {
	if dnsServer != nil {
		return dnsServer
	}
	var err error
	redisClient := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redisConfig.addr"),
		DB:   0,
	})
	dnsServer = &DNSServer{
		RedisClient: redisClient,
	}
	blockedIPs := viper.GetStringSlice("blocked_ip")
	for index := range blockedIPs {
		dnsServer.BlockedIP.Store(blockedIPs[index], struct{}{})
	}
	dnsServer.socket, err = net.ListenUDP(UDP_NETWORK, &net.UDPAddr{
		IP:   net.ParseIP(viper.GetString("dns_relay.client_ip")),
		Port: DNS_PORT,
	})
	if err != nil {
		log.Printf("network error or config error：%v", err)
		return dnsServer
	}

	log.Printf("There are %v ips blocked in total\n", len(blockedIPs))
	for index := range blockedIPs {
		log.Printf("Blocked ip：%v", blockedIPs[index])
	}
	return dnsServer
}

func (dnsServer *DNSServer) Serve() {
	for {
		data := make([]byte, 1024)
		read, remoteAddr, err := dnsServer.socket.ReadFromUDP(data)
		log.Println("-----------------------------------------------------")
		if err != nil {
			log.Println("Receive data error", err)
			continue
		}
		data = data[:read]
		log.Printf("Receive data success, length: %v, remote addr:%v, data: %v", read, remoteAddr, data)
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
