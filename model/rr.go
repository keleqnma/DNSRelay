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
	Name int16

	/* TYPE（2字节） */
	Type int16

	/* CLASS（2字节） */
	Class int16

	/* TTL（4字节） */
	TTL int

	/* RDLENGTH（2字节） */
	RDLength int16

	/* RDATA IPv4为4字节*/
	RData string
}

func NewDNSRR(Name int16, Type int16, Class int16, TTL int, RDLength int16, RData string) *DNSRR {
	return &DNSRR{
		Name:     Name,
		Type:     Type,
		Class:    Class,
		TTL:      TTL,
		RDLength: RDLength,
		RData:    RData,
	}
}

func (dnsRR *DNSRR) Pack() (data []byte, err error) {
	data = common.Pack(dnsRR.Name, dnsRR.Type, dnsRR.Class)
	data = append(data, common.IntToBytes4(dnsRR.TTL)...)
	data = append(data, common.Pack(dnsRR.RDLength)...)
	ipv4Data, err := common.Ipv4ToBytes(dnsRR.RData)
	data = append(data, ipv4Data...)
	return
}
