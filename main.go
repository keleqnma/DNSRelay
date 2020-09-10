package main

import (
	"fmt"

	"dnsrelay.com/v1/common"
	"dnsrelay.com/v1/model"
)

func main() {
	common.LoadConfig("conf")

	dnsHeader := model.NewDNSHeader(1, 3, 4, 5, 5, 6)
	dnsHeaderData := dnsHeader.PackDNSHeader()
	dnsHeaderProcess, _ := model.UnPackDNSHeader(dnsHeaderData)
	fmt.Println(dnsHeader, dnsHeaderData, dnsHeaderProcess, len(dnsHeaderData))

	dnsQuery := model.NewDNSQuery("test.com", 3, 4)
	dnsQueryData := dnsQuery.PackDNSQuery()
	dnsQueryProcess, _ := model.UnPackDNSQuery(dnsQueryData)
	fmt.Println(dnsQuery, dnsQueryData, dnsQueryProcess, len(dnsQueryData))

	dnsRR := model.NewDNSRR(12, 13, 12, 123, 11)
	dnsRRData := dnsQuery.PackDNSQuery()
	fmt.Println(dnsRR, dnsRRData, len(dnsRRData))
	// dnsServer := server.GetDNSServer()
	// dnsServer.Serve()
}
