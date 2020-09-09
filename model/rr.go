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
	Name int8

	/* TYPE（2字节） */
	Type int8

	/* CLASS（2字节） */
	Class int8

	/* TTL（4字节） */
	TTL int

	/* RDLENGTH（2字节） */
	RDLength int8

	/* RDATA IPv4为4字节*/
	RData string
}

func NewDNSRR(Name int8, Type int8, Class int8, TTL int, RDLength int8) *DNSRR {
	return &DNSRR{
		Name:     Name,
		Type:     Type,
		Class:    Class,
		TTL:      TTL,
		RDLength: RDLength,
	}
}

func (dnsRR *DNSRR) Pack() (data []byte) {
	data = common.Pack(dnsRR.Name, dnsRR.Type, dnsRR.Class)
	data = append(data, common.IntToBytes(dnsRR.TTL)...)
	data = append(data, common.Pack(dnsRR.RDLength)...)
	data = append(data, []byte(dnsRR.RData)...)
	return
}
