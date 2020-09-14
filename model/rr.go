package model

import "dnsrelay.com/v1/common"

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
type DNSRR struct {
	/* NAME (2字节) */
	Name int

	/* TYPE（2字节） */
	Type int

	/* CLASS（2字节） */
	Class int

	/* TTL（4字节） */
	TTL int

	/* RDLENGTH（2字节） */
	RDLength int

	/* RDATA IPv4为4字节*/
	RData string
}

func NewDNSRR(Name int, Type int, Class int, TTL int, RDLength int) *DNSRR {
	return &DNSRR{
		Name:     Name,
		Type:     Type,
		Class:    Class,
		TTL:      TTL,
		RDLength: RDLength,
	}
}

func (dnsRR *DNSRR) Pack() (data []byte, err error) {
	data = common.Pack(dnsRR.Name, dnsRR.Type, dnsRR.Class)
	data = append(data, common.IntToBytes4(dnsRR.TTL)...)
	data = append(data, common.Pack(dnsRR.RDLength)...)
	if len(dnsRR.RData) != 0 {
		ipv4Data, _ := common.Ipv4ToBytes(dnsRR.RData)
		data = append(data, ipv4Data...)
	}
	return
}
