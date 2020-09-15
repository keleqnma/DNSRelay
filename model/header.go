package model

import (
	"dnsrelay.com/v1/common"
)

const (
	HEADER_LENGTH   = 12
	HEADER_PACK_NUM = 6
	FAIL_FLAG       = 33155
	SUCCESS_FLAG    = 33152
)

/**
* DNS Header
   0  1  2  3  4  5  6  7  0  1  2  3  4  5  6  7
 +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
 |                      ID                       |
 +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
 |QR|  opcode   |AA|TC|RD|RA|   Z    |   RCODE   |
 +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
 |                    QDCOUNT                    |
 +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
 |                    ANCOUNT                    |
 +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
 |                    NSCOUNT                    |
 +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
 |                    ARCOUNT                    |
 +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type DNSHeader struct {
	/* 会话标识（2字节）*/
	ID int

	/* Flags（2字节）*/
	Flags int

	/* QDCOUNT（2字节）报文请求段中的问题记录数 */
	QDCount int

	/* ANCOUNT（2字节）报文回答段中的回答记录数*/
	ANCount int

	/* NSCOUNT（2字节）报文授权段中的授权记录数*/
	NSCount int

	/* ARCOUNT（2字节）报文附加段中的附加记录数*/
	ARCount int
}

func NewDNSHeader(ID int, Flags int, QDCount int, ANCount int, NSCount int, ARCount int) (dnsHeader *DNSHeader) {
	dnsHeader = &DNSHeader{}
	dnsHeader.ID, dnsHeader.Flags, dnsHeader.QDCount, dnsHeader.ANCount, dnsHeader.NSCount, dnsHeader.ARCount = ID, Flags, QDCount, ANCount, NSCount, ARCount
	return
}

func UnPackDNSHeader(data []byte) (dnsHeader *DNSHeader) {
	nums := common.UnPack(data)
	dnsHeader = NewDNSHeader(nums[0], nums[1], nums[2], nums[3], nums[4], nums[5])
	return dnsHeader
}

func (dnsHeader *DNSHeader) PackDNSHeader() (data []byte) {
	data = common.Pack(dnsHeader.ID, dnsHeader.Flags, dnsHeader.QDCount, dnsHeader.ANCount, dnsHeader.NSCount, dnsHeader.ARCount)
	return
}
