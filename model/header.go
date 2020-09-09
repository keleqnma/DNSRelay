package model

import (
	"errors"

	"dnsrelay.com/v1/common"
)

const (
	HEADER_LENGTH   = 12
	HEADER_PACK_NUM = 6
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
	ID int8

	/* Flags（2字节）*/
	Flags int8

	/* QDCOUNT（2字节）*/
	QDCount int8

	/* ANCOUNT（2字节）*/
	ANCount int8

	/* NSCOUNT（2字节）*/
	NSCount int8

	/* ARCOUNT（2字节）*/
	ARCount int8
}

func NewDNSHeader(ID int8, Flags int8, QDCount int8, ANCount int8, NSCount int8, ARCount int8) (dnsHeader *DNSHeader) {
	dnsHeader = &DNSHeader{}
	dnsHeader.ID, dnsHeader.Flags, dnsHeader.QDCount, dnsHeader.ANCount, dnsHeader.NSCount, dnsHeader.ARCount = ID, Flags, QDCount, ANCount, NSCount, ARCount
	return
}

func UnPackDNSHeader(data []byte) (dnsHeader *DNSHeader, err error) {
	nums := common.UnPack(data)
	if len(nums) != HEADER_PACK_NUM {
		err = errors.New("dns header 解析失败")
	}
	dnsHeader = NewDNSHeader(nums[0], nums[1], nums[2], nums[3], nums[4], nums[5])
	return dnsHeader, err
}

func (dnsHeader *DNSHeader) PackDNSHeader() (data []byte) {
	data = common.Pack(dnsHeader.ID, dnsHeader.Flags, dnsHeader.QDCount, dnsHeader.ANCount, dnsHeader.NSCount, dnsHeader.ARCount)
	return
}
