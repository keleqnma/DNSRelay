package model

import (
	"errors"

	"dnsrelay.com/v1/common"
)

const (
	HOST_QUERY_TYPE = 1
	QUERY_PACK_NUM  = 2
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
type DNSQuestion struct {
	/* QNAME 8bit为单位表示的查询名(广泛的说就是：域名) */
	QName string

	/* QTYPE（2字节）查询的协议类型*/
	QType int

	/* QCLASS（2字节）查询的类,比如，IN代表Internet */
	QClass int
}

func NewDNSQuestion(QName string, QType int, QClass int) (dnsQuestion *DNSQuestion) {
	dnsQuestion = &DNSQuestion{}
	dnsQuestion.QName, dnsQuestion.QType, dnsQuestion.QClass = QName, QType, QClass
	return dnsQuestion
}

func UnPackDNSQuestion(data []byte) (dnsQuestion *DNSQuestion, err error) {
	dnsQuestion = &DNSQuestion{}
	nameLength := 0
	nameLength, dnsQuestion.QName = common.BytesToDomain(data)
	nums := common.UnPack(data[nameLength:])
	if len(nums) < QUERY_PACK_NUM {
		err = errors.New("dns query 解析失败")
	}
	dnsQuestion.QType, dnsQuestion.QClass = nums[0], nums[1]
	return dnsQuestion, err
}

func (dnsQuestion *DNSQuestion) PackDNSQuestion() (data []byte) {
	data = append(data, common.DomainToBytes(dnsQuestion.QName)...)
	data = append(data, common.Pack(dnsQuestion.QType, dnsQuestion.QClass)...)
	return
}
