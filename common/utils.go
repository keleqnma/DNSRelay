package common

import (
	"bytes"
	"encoding/binary"
)

func Pack(args ...int8) []byte {
	ret := []byte{}
	for i := range args {
		ret = append(ret, Int8ToBytes2(args[i])...)
	}
	return ret
}

func UnPack(data []byte) []int8 {
	ret := []int8{}
	remainLen := len(data)
	for remainLen >= 2 {
		ret = append(ret, BytesToInt8(data[0:2]))
		remainLen -= 2
		data = data[2:]
	}
	return ret
}

func Int8ToBytes2(num int8) []byte {
	m := int32(num)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, m) //nolint: errcheck

	gbyte := bytesBuffer.Bytes()
	k := 2
	x := len(gbyte)
	nb := make([]byte, k)
	for i := 0; i < k; i++ {
		nb[i] = gbyte[x-i-1]
	}
	return nb
}

func IntToBytes(num int) []byte {
	m := int32(num)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, m) //nolint: errcheck

	gbyte := bytesBuffer.Bytes()
	k := 4
	x := len(gbyte)
	nb := make([]byte, k)
	for i := 0; i < k; i++ {
		nb[i] = gbyte[x-i-1]
	}
	return nb
}

func BytesToInt8(b []byte) int8 {
	var xx []byte
	if len(b) == 2 {
		xx = []byte{b[0], b[1], 0, 0}
	} else {
		xx = b
	}

	m := len(xx)
	nb := make([]byte, 4)
	for i := 0; i < 4; i++ {
		nb[i] = xx[m-i-1]
	}
	bytesBuffer := bytes.NewBuffer(nb)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x) //nolint: errcheck

	return int8(x)
}
