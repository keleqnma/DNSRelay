package server

import (
	"log"
	"net"

	"dnsrelay.com/v1/model"
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
	if parserServer.headerReceived, err = model.UnPackDNSHeader(parserServer.dataReceived); err != nil {
		log.Printf("Header unpacked failed : %v\n", err)
		return
	}

	if parserServer.headerReceived.QDCount <= 0 {
		log.Printf("Header received %v length error\n", parserServer.headerReceived)
		return
	}

	if parserServer.queryReceived, err = model.UnPackDNSQuery(parserServer.dataReceived[model.HEADER_LENGTH:]); err != nil {
		log.Printf("Query unpacked failed : %v\n", err)
		return
	}

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
	if _, ok := dnsServer.BlockedIP[ipSearch]; ok {
		flag = model.ERROR_FLAG
	} else {
		flag = model.SUCCESS_FLAG
	}

	dnsHeaderResp = model.NewDNSHeader(parserServer.headerReceived.ID, flag, parserServer.headerReceived.QDCount, AN_COUNT_INIT, NS_COUNT_INIT, AR_COUNT_INIT)
	respData = append(respData, dnsHeaderResp.PackDNSHeader()...)

	dnsQueryResp = parserServer.queryReceived
	respData = append(respData, dnsQueryResp.PackDNSQuery()...)

	dnsRRResp = model.NewDNSRR(RR_NAME_INIT, parserServer.queryReceived.QType, parserServer.queryReceived.QClass, RR_TTL_INIT, RR_RD_LEN_INIT, ipSearch)
	dnsRRRespData, err := dnsRRResp.Pack()
	if err != nil {
		log.Printf("Invalid dnsRRResp format;%v", err)
	}
	respData = append(respData, dnsRRRespData...)

	dnsRRNameServer = model.NewDNSRR(RR_NAME_INIT, RR_AUTHOR_TYPE_INIT, parserServer.queryReceived.QClass, RR_TTL_INIT, RR_AUTHOR_RD_LEN_INIT, "")
	dnsRRNameServerData, err := dnsRRNameServer.Pack()
	if err != nil {
		log.Printf("Invalid dnsRRNameServer format;%v", err)
	}
	respData = append(respData, dnsRRNameServerData...)

	log.Printf("Local search responce data：%v", respData)

	code, err := GetDNSServer().socket.WriteToUDP(respData, parserServer.clientAddr)
	if err != nil {
		log.Printf("Local server write error:%v, code %v \n", err, code)
	}

	log.Printf("Search local server done. domian：%s，ip searched：%s\n", parserServer.queryReceived.QName, ipSearch)
	return
}

func (parserServer *ParserServer) searchInternet() {
	var (
		code      int
		err       error
		dataTrans []byte
	)
	parserServer.transSocket, err = net.ListenUDP(UDP_NETWORK, TRANS_ADDR)
	if err != nil {
		log.Panicf("Net server config error：%v", err)
		return
	}
	code, err = parserServer.transSocket.WriteToUDP(parserServer.dataReceived, parserServer.clientAddr)
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
	log.Printf("Search net server done. domian：%s\n", parserServer.queryReceived.QName)
}
