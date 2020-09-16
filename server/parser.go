package server

import (
	"context"
	"log"
	"net"
	"strconv"

	"dnsrelay.com/v1/common"
	"dnsrelay.com/v1/model"
	"github.com/spf13/viper"
)

type ParserServer struct {
	clientAddr *net.UDPAddr

	dataReceived   []byte
	headerReceived *model.DNSHeader
	queryReceived  *model.DNSQuestion
}

var ctx = context.Background()

func GetParserServer(data []byte, addr *net.UDPAddr) (*ParserServer, error) {
	var err error
	parserServer := &ParserServer{}
	parserServer.dataReceived = data
	parserServer.clientAddr = addr
	return parserServer, err
}

func (parserServer *ParserServer) parse() (err error) {
	parserServer.headerReceived = model.UnPackDNSHeader(parserServer.dataReceived[:model.HEADER_LENGTH])
	defer log.Println("-----------------------------------------------------")
	if parserServer.headerReceived.QDCount <= 0 {
		log.Printf("Header received %v length error\n", parserServer.headerReceived)
		return
	}

	if parserServer.queryReceived, err = model.UnPackDNSQuestion(parserServer.dataReceived[model.HEADER_LENGTH:]); err != nil {
		log.Printf("Query unpacked failed : %v\n", err)
		return
	}

	if parserServer.queryReceived.QType == model.HOST_QUERY_TYPE && parserServer.searchLocal() {
		return
	}

	parserServer.searchInternet()

	return
}

func (parserServer *ParserServer) searchLocal() (ok bool) {
	log.Printf("Search local server for domian：%s\n", parserServer.queryReceived.QName)
	searchKey := DNS_PROXY_REDIS_SPACE + parserServer.queryReceived.QName
	ips, err := dnsServer.RedisClient.SMembers(ctx, searchKey).Result()
	if err != nil {
		log.Printf("redis err: %v", err)
		return
	}
	if len(ips) == 0 {
		log.Printf("Search local server for domian：%s not found\n", parserServer.queryReceived.QName)
		return
	}

	flag := model.SUCCESS_FLAG
	for index := range ips {
		_, ok = dnsServer.BlockedIP.Load(ips[index])
		if ok {
			flag = model.FAIL_FLAG
			log.Printf("Search net server search domian：%s，ip blocked：%s\n", parserServer.queryReceived.QName, ips[index])
			break
		}
	}
	if flag == model.SUCCESS_FLAG {
		log.Printf("Search local server search domian：%s，ip searched：%s\n", parserServer.queryReceived.QName, ips)
	}

	var respData []byte

	respData = append(respData, model.NewDNSHeader(parserServer.headerReceived.ID, flag, parserServer.headerReceived.QDCount, len(ips), NS_COUNT_INIT, AR_COUNT_INIT).PackDNSHeader()...)

	respData = append(respData, parserServer.dataReceived[model.HEADER_LENGTH:]...)

	for index := range ips {
		respData = append(respData, model.NewDNSAnswer(ANSWER_NAME_INIT, parserServer.queryReceived.QType, parserServer.queryReceived.QClass, ANSWER_TTL_INIT, ANSWER_RD_LEN_INIT, ips[index]).Pack()...)
	}

	length, err := GetDNSServer().socket.WriteToUDP(respData, parserServer.clientAddr)
	log.Printf("Local search server send length:%v, data：%v", length, respData)
	if err != nil {
		log.Printf("Local server write error:%v, length %v \n", err, length)
	}
	log.Printf("Search local server for domian：%s done\n", parserServer.queryReceived.QName)
	return true
}

func (parserServer *ParserServer) searchInternet() {
	var (
		length int
		err    error
	)
	dataTrans := make([]byte, 1024)
	log.Printf("Search net server for domian：%s\n", parserServer.queryReceived.QName)

	dstServer := &net.UDPAddr{
		IP:   net.ParseIP(viper.GetString("dns_relay.trans_ip")),
		Port: DNS_PORT,
	}
	srcServer := &net.UDPAddr{IP: net.IPv4zero, Port: 0}

	conn, err := net.DialUDP(UDP_NETWORK, srcServer, dstServer)
	if err != nil {
		log.Panicf("Net server Listen error：%v", err)
		return
	}

	length, err = conn.Write(parserServer.dataReceived)
	log.Printf("send data to transport server:%v , length:%v\n", dstServer, length)
	if err != nil {
		log.Printf("Net server write error:%v, length %v \n", err, length)
		return
	}

	length, err = conn.Read(dataTrans)
	if err != nil {
		log.Printf("Net server read error:%v, length %v \n", err, length)
		return
	}
	dataTrans = dataTrans[:length]
	if conn != nil {
		defer conn.Close()
	}

	dataSend := dataTrans

	if parserServer.queryReceived.QType == model.HOST_QUERY_TYPE {
		transHeader := model.UnPackDNSHeader(dataTrans[:model.HEADER_LENGTH])
		answerNums := transHeader.ANCount
		key := DNS_PROXY_REDIS_SPACE + parserServer.queryReceived.QName
		sendSize := len(parserServer.dataReceived)
		dataTrans = dataTrans[sendSize:]
		var values []string

		for num := 0; num < answerNums && len(dataTrans) >= model.IPV4_ANSWER_LEANGTH; num++ {
			dataTrans = dataTrans[(model.IPV4_ANSWER_LEANGTH - model.IPV4_RDATA_LENGTH):]
			var ip string
			for index := 0; index < model.IPV4_RDATA_LENGTH; index++ {
				ip += strconv.Itoa(int(dataTrans[index]))
				ip += "."
			}
			ip = ip[:len(ip)-1]
			if _, ok := dnsServer.BlockedIP.Load(ip); ok {
				log.Printf("Search net server done. domian：%s，ip blocked：%s\n", parserServer.queryReceived.QName, ip)
				failFlags := common.IntToBytes2(model.FAIL_FLAG)
				dataSend[2], dataSend[3] = failFlags[0], failFlags[1]
				break
			}
			values = append(values, ip)
			dataTrans = dataTrans[model.IPV4_RDATA_LENGTH:]
		}

		if len(values) > 0 {
			if dnsServer.RedisClient.SAdd(ctx, key, values).Err() != nil {
				log.Printf("write to local database error: %v.", err)
			} else {
				log.Printf("write to local database success, domain:%v, ips:%v.", parserServer.queryReceived.QName, values)
			}
		}
	}

	length, err = dnsServer.socket.WriteToUDP(dataSend, parserServer.clientAddr)
	if err != nil {
		log.Printf("Net server write error:%v, length %v \n", err, length)
	}
	log.Printf("Search net server send length：%v, data: %v \n", length, dataSend)
	log.Printf("Search net server for domian：%s done\n", parserServer.queryReceived.QName)
}
