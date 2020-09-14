package server

import "net"

const (
	UDP_NETWORK = "udp4"
	TRANS_PORT  = 53

	AN_COUNT_INIT = 1
	NS_COUNT_INIT = 1
	AR_COUNT_INIT = 0

	RR_NAME_INIT   = 0x00c
	RR_TTL_INIT    = 3600 * 24
	RR_RD_LEN_INIT = 4

	RR_AUTHOR_RD_LEN_INIT = 0
	RR_AUTHOR_TYPE_INIT   = 6
)

var (
	TRANS_ADDR = &net.UDPAddr{
		IP:   net.IPv4(114, 114, 114, 114),
		Port: TRANS_PORT,
	}
)
