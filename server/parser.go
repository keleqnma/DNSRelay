package server

import (
	"log"
	"net"

	"dnsrelay.com/v1/model"
	"github.com/spf13/viper"
)

type ParserServer struct {
	clientAddr *net.UDPAddr

	dataReceived   []byte
	headerReceived *model.DNSHeader
	queryReceived  *model.DNSQuestion
}

func GetParserServer(data []byte, addr *net.UDPAddr) (*ParserServer, error) {
	var err error
	parserServer := &ParserServer{}
	parserServer.dataReceived = data
	parserServer.clientAddr = addr
	return parserServer, err
}

func (parserServer *ParserServer) parse() (err error) {
	parserServer.headerReceived = model.UnPackDNSHeader(parserServer.dataReceived[:model.HEADER_LENGTH])
	if parserServer.headerReceived.QDCount <= 0 {
		log.Printf("Header received %v length error\n", parserServer.headerReceived)
		return
	}

	if parserServer.queryReceived, err = model.UnPackDNSQuestion(parserServer.dataReceived[model.HEADER_LENGTH:]); err != nil {
		log.Printf("Query unpacked failed : %v\n", err)
		return
	}

	if ok := parserServer.searchLocal(); ok {
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
	ipSearch, ok := dnsServer.DomainMap[parserServer.queryReceived.QName]
	if !ok || parserServer.queryReceived.QType != model.HOST_QUERY_TYPE {
		log.Printf("Search local server for domian：%s not found.\n", parserServer.queryReceived.QName)
		ok = false
		return
	}

	var flag int
	if _, ok := dnsServer.BlockedIP[ipSearch]; ok {
		flag = model.FAIL_FLAG
		dnsHeaderResp = model.NewDNSHeader(parserServer.headerReceived.ID, flag, parserServer.headerReceived.QDCount, AN_FAIL_COUNT_INIT, NS_COUNT_INIT, AR_COUNT_INIT)
		log.Printf("Search local server done. domian：%s，ip blocked：%s\n", parserServer.queryReceived.QName, ipSearch)
	} else {
		flag = model.SUCCESS_FLAG
		dnsHeaderResp = model.NewDNSHeader(parserServer.headerReceived.ID, flag, parserServer.headerReceived.QDCount, AN_SUC_COUNT_INIT, NS_COUNT_INIT, AR_COUNT_INIT)
		log.Printf("Search local server done. domian：%s，ip searched：%s\n", parserServer.queryReceived.QName, ipSearch)
	}
	respData = append(respData, dnsHeaderResp.PackDNSHeader()...)

	dnsQuestionResp = parserServer.queryReceived
	respData = append(respData, dnsQuestionResp.PackDNSQuestion()...)

	dnsAnswerResp = model.NewDNSAnswer(ANSWER_NAME_INIT, parserServer.queryReceived.QType, parserServer.queryReceived.QClass, ANSWER_TTL_INIT, ANSWER_RD_LEN_INIT, ipSearch)
	respData = append(respData, dnsAnswerResp.Pack()...)

	code, err := GetDNSServer().socket.WriteToUDP(respData, parserServer.clientAddr)
	log.Printf("Local search server send code:%v, data：%v", code, respData)
	if err != nil {
		log.Printf("Local server write error:%v, code %v \n", err, code)
	}
	return true
}

func (parserServer *ParserServer) searchInternet() {
	var (
		code int
		err  error
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

	code, err = conn.Write(parserServer.dataReceived)
	log.Printf("send data to transport server:%v , code:%v\n", dstServer, code)
	if err != nil {
		log.Printf("Net server write error:%v, code %v \n", err, code)
		return
	}

	code, err = conn.Read(dataTrans)
	if err != nil {
		log.Printf("Net server read error:%v, code %v \n", err, code)
		return
	}
	dataTrans = dataTrans[:code]

	if conn != nil {
		defer conn.Close()
	}

	code, err = dnsServer.socket.WriteToUDP(dataTrans, parserServer.clientAddr)
	if err != nil {
		log.Printf("Net server write error:%v, code %v \n", err, code)
	}
	log.Printf("Search net server send code：%v, data: %v \n", code, dataTrans)
	log.Printf("Search net server for domian：%s done\n", parserServer.queryReceived.QName)
}
