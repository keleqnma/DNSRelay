package model

import (
	"errors"

	"dnsrelay.com/v1/common"
)

const (
	HOST_QUERY_TYPE = 1
	QUERY_PACK_NUM  = 2
	dotByte         = 32
)

/**
 * Query 查询字段
	0  1  2  3  4  5  6  7  0  1  2  3  4  5  6  7
  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
  |                     ...                       |
  |                    QNAME                      |
  |                     ...                       |
  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
  |                    QTYPE                      |
  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
  |                    QCLASS                     |
  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type DNSQuery struct {
	/* QNAME 8bit为单位表示的查询名(广泛的说就是：域名) */
	QName string

	/* QTYPE（2字节） */
	QType int8

	/* QCLASS（2字节） */
	QClass int8
}

func NewDNSQuery(QName string, QType int8, QClass int8) (dnsQuery *DNSQuery) {
	dnsQuery = &DNSQuery{}
	dnsQuery.QName, dnsQuery.QType, dnsQuery.QClass = QName, QType, QClass
	return dnsQuery
}

func UnPackDNSQuery(data []byte) (dnsQuery *DNSQuery, err error) {
	dnsQuery = &DNSQuery{}
	nameLength := 0
	nameLength, dnsQuery.QName = common.BytesToDomain(data)
	nums := common.UnPack(data[nameLength:])
	if len(nums) != QUERY_PACK_NUM {
		err = errors.New("dns header 解析失败")
	}
	dnsQuery.QType, dnsQuery.QClass = nums[0], nums[1]
	return dnsQuery, err
}

func (dnsQuery *DNSQuery) PackDNSQuery() (data []byte) {
	data = append(data, common.DomainToBytes(dnsQuery.QName)...)
	data = append(data, common.Pack(dnsQuery.QType, dnsQuery.QClass)...)
	return
}
