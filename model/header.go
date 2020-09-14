package model

import (
	"errors"

	"dnsrelay.com/v1/common"
)

const (
	HEADER_LENGTH   = 12
	HEADER_PACK_NUM = 6
	ERROR_FLAG      = 33155
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
	ID int16

	/* Flags（2字节）*/
	Flags int

	/* QDCOUNT（2字节）*/
	QDCount int16

	/* ANCOUNT（2字节）*/
	ANCount int16

	/* NSCOUNT（2字节）*/
	NSCount int16

	/* ARCOUNT（2字节）*/
	ARCount int16
}

func NewDNSHeader(ID int16, Flags int, QDCount int16, ANCount int16, NSCount int16, ARCount int16) (dnsHeader *DNSHeader) {
	dnsHeader = &DNSHeader{}
	dnsHeader.ID, dnsHeader.Flags, dnsHeader.QDCount, dnsHeader.ANCount, dnsHeader.NSCount, dnsHeader.ARCount = ID, Flags, QDCount, ANCount, NSCount, ARCount
	return
}

func UnPackDNSHeader(data []byte) (dnsHeader *DNSHeader, err error) {
	id := common.BytesToInt16(data[:2])
	flag := common.BytesToInt(data[2:4])
	nums := common.UnPack(data[4:])
	if len(nums) != HEADER_PACK_NUM {
		err = errors.New("dns header 解析失败")
	}
	dnsHeader = NewDNSHeader(id, flag, nums[0], nums[1], nums[2], nums[3])
	return dnsHeader, err
}

func (dnsHeader *DNSHeader) PackDNSHeader() (data []byte) {
	data = common.Pack(dnsHeader.ID)
	data = append(data, common.IntToBytes2(dnsHeader.Flags)...)
	data = append(data, common.Pack(dnsHeader.QDCount, dnsHeader.ANCount, dnsHeader.NSCount, dnsHeader.ARCount)...)
	return
}
