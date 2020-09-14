package server

import (
	"log"
	"net"

	"dnsrelay.com/v1/model"
	"github.com/spf13/viper"
)

type ParserServer struct {
	transSocket *net.UDPConn
	clientAddr  *net.UDPAddr

	dataReceived   []byte
	headerReceived *model.DNSHeader
	queryReceived  *model.DNSQuery
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
	log.Printf("receive header:%v, unpacked header:%v\n", parserServer.dataReceived[:model.HEADER_LENGTH], parserServer.headerReceived)
	if parserServer.headerReceived.QDCount <= 0 {
		log.Printf("Header received %v length error\n", parserServer.headerReceived)
		return
	}

	if parserServer.queryReceived, err = model.UnPackDNSQuery(parserServer.dataReceived[model.HEADER_LENGTH:]); err != nil {
		log.Printf("Query unpacked failed : %v\n", err)
		return
	}
	log.Printf("receive query:%v, unpacked query:%v\n", parserServer.dataReceived[model.HEADER_LENGTH:], parserServer.queryReceived)

	if ok := parserServer.searchLocal(); ok {
		return
	}

	parserServer.searchInternet()

	if parserServer.transSocket != nil {
		defer parserServer.transSocket.Close()
	}
	return
}

func (parserServer *ParserServer) searchLocal() (ok bool) {
	var (
		respData        []byte
		dnsHeaderResp   *model.DNSHeader
		dnsQueryResp    *model.DNSQuery
		dnsRRResp       *model.DNSRR
		dnsRRNameServer *model.DNSRR
	)
	ipSearch, ok := dnsServer.DomainMap[parserServer.queryReceived.QName]
	if !ok || parserServer.queryReceived.QType != model.IPV4_QUERY_TYPE {
		ok = false
		return
	}

	var flag int
	var anCount int
	if _, ok := dnsServer.BlockedIP[ipSearch]; ok {
		flag = model.ERROR_FLAG
		anCount = AN_FAIL_COUNT_INIT
	} else {
		flag = model.SUCCESS_FLAG
		anCount = AN_SUC_COUNT_INIT
	}

	dnsHeaderResp = model.NewDNSHeader(parserServer.headerReceived.ID, flag, parserServer.headerReceived.QDCount, anCount, NS_COUNT_INIT, AR_COUNT_INIT)
	respData = append(respData, dnsHeaderResp.PackDNSHeader()...)

	dnsQueryResp = parserServer.queryReceived
	respData = append(respData, dnsQueryResp.PackDNSQuery()...)

	if flag == model.SUCCESS_FLAG {
		dnsRRResp = model.NewDNSRR(RR_NAME_INIT, parserServer.queryReceived.QType, parserServer.queryReceived.QClass, RR_TTL_INIT, RR_RD_LEN_INIT)
		dnsRRResp.RData = ipSearch
		dnsRRRespData, err := dnsRRResp.Pack()
		if err != nil {
			log.Printf("Invalid dnsRRResp format;%v", err)
		}
		respData = append(respData, dnsRRRespData...)

		dnsRRNameServer = model.NewDNSRR(RR_NAME_INIT, RR_AUTHOR_TYPE_INIT, parserServer.queryReceived.QClass, RR_TTL_INIT, RR_AUTHOR_RD_LEN_INIT)
		dnsRRNameServerData, _ := dnsRRNameServer.Pack()
		respData = append(respData, dnsRRNameServerData...)

		log.Printf("Search local server done. domian：%s，ip searched：%s\n", parserServer.queryReceived.QName, ipSearch)
	}

	log.Printf("Local search responce data：%v", respData)
	log.Println(dnsHeaderResp, dnsQueryResp, dnsRRResp, dnsRRNameServer)
	code, err := GetDNSServer().socket.WriteToUDP(respData, parserServer.clientAddr)
	if err != nil {
		log.Printf("Local server write error:%v, code %v \n", err, code)
	}
	return
}

func (parserServer *ParserServer) searchInternet() {
	var (
		code int
		err  error
	)
	dataTrans := make([]byte, 1024)
	log.Printf("Search net server for domian：%s\n", parserServer.queryReceived.QName)

	parserServer.transSocket, err = net.ListenUDP(UDP_NETWORK, &net.UDPAddr{Port: 0})
	if err != nil {
		log.Panicf("Net server Listen error：%v", err)
		return
	}

	code, err = parserServer.transSocket.WriteToUDP(parserServer.dataReceived, &net.UDPAddr{
		IP:   net.ParseIP(viper.GetString("dns_relay.trans_ip")),
		Port: 53,
	})
	log.Printf("send data to transport server, code:%v\n", code)
	if err != nil {
		log.Printf("Net server write error:%v, code %v \n", err, code)
		return
	}

	code, _, err = parserServer.transSocket.ReadFromUDP(dataTrans)
	if err != nil {
		log.Printf("Net server read error:%v, code %v \n", err, code)
		return
	}

	code, err = dnsServer.socket.WriteToUDP(dataTrans, parserServer.clientAddr)
	if err != nil {
		log.Printf("Net server write error:%v, code %v \n", err, code)
	}
	log.Printf("Search net server for domian：%s done\n", parserServer.queryReceived.QName)
}
