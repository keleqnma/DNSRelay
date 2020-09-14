package main

import (
	"dnsrelay.com/v1/common"
	"dnsrelay.com/v1/server"
)

func main() {
	common.LoadConfig("conf")

	// dnsHeader := model.NewDNSHeader(1, 3, 4, 5, 5, 6)
	// dnsHeaderData := dnsHeader.PackDNSHeader()
	// dnsHeaderProcess, _ := model.UnPackDNSHeader(dnsHeaderData)
	// fmt.Println(dnsHeader, dnsHeaderData, dnsHeaderProcess, len(dnsHeaderData))

	// dnsQuery := model.NewDNSQuery("test.chenyuqi.yuqi.com", 3, 4)
	// dnsQueryData := dnsQuery.PackDNSQuery()
	// dnsQueryProcess, _ := model.UnPackDNSQuery(dnsQueryData)
	// fmt.Println(dnsQuery, dnsQueryData, dnsQueryProcess, len(dnsQueryData))

	// dnsRR := model.NewDNSRR(12, 13, 12, 123, 11, "255.234.123.112")
	// dnsRRData, err := dnsRR.Pack()
	// fmt.Println(dnsRR, dnsRRData, len(dnsRRData), err)

	dnsServer := server.GetDNSServer()
	dnsServer.Serve()
}
