package server

import (
	"context"
	"log"
	"net"
	"strconv"

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
	var (
		respData        []byte
		dnsHeaderResp   *model.DNSHeader
		dnsQuestionResp *model.DNSQuestion
		dnsAnswerResp   *model.DNSAnswer
	)
	log.Printf("Search local server for domian：%s\n", parserServer.queryReceived.QName)
	searchKey := DNS_PROXY_REDIS_SPACE + parserServer.queryReceived.QName
	ipSearchs, err := dnsServer.RedisClient.SMembers(ctx, searchKey).Result()
	if err != nil {
		return
	}
	if len(ipSearchs) == 0 {
		log.Printf("Search local server for domian：%s not found, err:%v.\n", parserServer.queryReceived.QName, err)
		return
	}

	for index := range ipSearchs {
		_, ok = dnsServer.BlockedIP.Load(ipSearchs[index])
		if ok {
			break
		}
	}

	var flag int
	if ok {
		flag = model.FAIL_FLAG
		dnsHeaderResp = model.NewDNSHeader(parserServer.headerReceived.ID, flag, parserServer.headerReceived.QDCount, AN_FAIL_COUNT_INIT, NS_COUNT_INIT, AR_COUNT_INIT)
		log.Printf("Search local server done. domian：%s，ip blocked：%s\n", parserServer.queryReceived.QName, ipSearchs)
	} else {
		flag = model.SUCCESS_FLAG
		dnsHeaderResp = model.NewDNSHeader(parserServer.headerReceived.ID, flag, parserServer.headerReceived.QDCount, len(ipSearchs), NS_COUNT_INIT, AR_COUNT_INIT)
		log.Printf("Search local server done. domian：%s，ip searched：%s\n", parserServer.queryReceived.QName, ipSearchs)
	}
	respData = append(respData, dnsHeaderResp.PackDNSHeader()...)

	dnsQuestionResp = parserServer.queryReceived
	respData = append(respData, dnsQuestionResp.PackDNSQuestion()...)

	for index := range ipSearchs {
		dnsAnswerResp = model.NewDNSAnswer(ANSWER_NAME_INIT, parserServer.queryReceived.QType, parserServer.queryReceived.QClass, ANSWER_TTL_INIT, ANSWER_RD_LEN_INIT, ipSearchs[index])
		respData = append(respData, dnsAnswerResp.Pack()...)
	}

	length, err := GetDNSServer().socket.WriteToUDP(respData, parserServer.clientAddr)
	log.Printf("Local search server send length:%v, data：%v", length, respData)
	if err != nil {
		log.Printf("Local server write error:%v, length %v \n", err, length)
	}
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
		key := DNS_PROXY_REDIS_SPACE + parserServer.queryReceived.QName
		sendSize := len(parserServer.dataReceived)
		dataTrans = dataTrans[sendSize:]
		var values []string

		for len(dataTrans) >= model.IPV4_ANSWER_LEANGTH {
			dataTrans = dataTrans[(model.IPV4_ANSWER_LEANGTH - model.IPV4_RDATA_LENGTH):]
			var ip string
			for index := 0; index < model.IPV4_RDATA_LENGTH; index++ {
				ip += strconv.Itoa(int(dataTrans[index]))
				ip += "."
			}
			values = append(values, ip[:len(ip)-1])
			dataTrans = dataTrans[model.IPV4_RDATA_LENGTH:]
		}
		if err := dnsServer.RedisClient.SAdd(ctx, key, values).Err(); err != nil {
			log.Printf("write to local database error: %v.", err)
		} else {
			log.Printf("write to local database success, domain:%v, ips:%v.", parserServer.queryReceived.QName, values)
		}
	}

	length, err = dnsServer.socket.WriteToUDP(dataSend, parserServer.clientAddr)
	if err != nil {
		log.Printf("Net server write error:%v, length %v \n", err, length)
	}
	log.Printf("Search net server send length：%v, data: %v \n", length, dataSend)
	log.Printf("Search net server for domian：%s done\n", parserServer.queryReceived.QName)
}
