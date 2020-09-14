package main

import (
	"dnsrelay.com/v1/common"
	"dnsrelay.com/v1/server"
)

func main() {
	common.LoadConfig("conf")

	dnsServer := server.GetDNSServer()
	dnsServer.Serve()
}
