package model

import "dnsrelay.com/v1/common"

const (
	IPV4_ANSWER_LEANGTH = 16
	IPV4_RDATA_LENGTH   = 4
)

/**
	   0  1  2  3  4  5  6  7  0  1  2  3  4  5  6  7
	  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	  |					   ... 						  |
	  |                    NAME                       |
	  |                    ...                        |
	  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	  |                    TYPE                       |
	  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	  |                    CLASS                      |
	  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	  |                    TTL                        |
      |                                               |
	  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	  |                    RDLENGTH                   |
	  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	  |                    ...                        |
	  |                    RDATA                      |
	  |                    ...                        |
	  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type DNSAnswer struct {
	/* NAME (2字节) 记录包含的域名*/
	Name int

	/* TYPE（2字节）DNS协议的类型 */
	Type int

	/* CLASS（2字节） RDATA的类*/
	Class int

	/* TTL（4字节） 表示资源记录可以缓存的时间。0代表只能被传输，但是不能被缓存。*/
	TTL int

	/* RDLENGTH（2字节） RDATA的长度*/
	RDLength int

	/* RDATA IPv4为4字节 不定长字符串来表示记录，格式根TYPE和CLASS有关。比如，TYPE是A，CLASS 是 IN，那么RDATA就是一个4个字节的ARPA网络地址。*/
	RData string
}

func NewDNSAnswer(Name int, Type int, Class int, TTL int, RDLength int, RData string) *DNSAnswer {
	return &DNSAnswer{
		Name:     Name,
		Type:     Type,
		Class:    Class,
		TTL:      TTL,
		RDLength: RDLength,
		RData:    RData,
	}
}

func (dnsAnswer *DNSAnswer) Pack() (data []byte) {
	data = common.Pack(dnsAnswer.Name, dnsAnswer.Type, dnsAnswer.Class)
	data = append(data, common.IntToBytes4(dnsAnswer.TTL)...)
	data = append(data, common.Pack(dnsAnswer.RDLength)...)
	if len(dnsAnswer.RData) != 0 {
		ipv4Data, _ := common.Ipv4ToBytes(dnsAnswer.RData)
		data = append(data, ipv4Data...)
	}
	return
}
