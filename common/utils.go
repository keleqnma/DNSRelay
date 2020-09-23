package common

import (
	"errors"
	"strconv"
	"strings"
)

const stopByte = 0x00

func Pack(args ...int) []byte {
	ret := []byte{}
	for i := range args {
		ret = append(ret, IntToBytes2(args[i])...)
	}
	return ret
}

func UnPack(data []byte) []int {
	ret := []int{}
	remainLen := len(data)
	for remainLen >= 2 {
		ret = append(ret, Bytes2ToInt(data[0:2]))
		remainLen -= 2
		data = data[2:]
	}
	return ret
}

func IntToBytes4(num int) []byte {
	m := int32(num)
	var res []byte
	res = append(res, byte((m>>24)&0xFF))
	res = append(res, byte((m>>16)&0xFF))
	res = append(res, byte((m>>8)&0xFF))
	res = append(res, byte((m)&0xFF))
	return res
}

func IntToBytes2(num int) []byte {
	m := int32(num)
	var res []byte
	res = append(res, byte((m>>8)&0xFF))
	res = append(res, byte((m)&0xFF))
	return res
}

func Bytes2ToInt(b []byte) int {
	byte1 := int(b[0]&0xff) << 8
	byte2 := int(b[1] & 0xff)
	return byte1 | byte2
}

func Bytes4ToInt(b []byte) int {
	byte1 := int(b[0]&0xff) << 24
	byte2 := int(b[1]&0xff) << 16
	byte3 := int(b[2]&0xff) << 8
	byte4 := int(b[3] & 0xff)
	return byte1 | byte2 | byte3 | byte4
}

func DomainToBytes(domain string) []byte {
	var res []byte
	nums := strings.Split(domain, ".")
	for index := range nums {
		res = append(res, byte(len(nums[index])))
		res = append(res, []byte(nums[index])...)
	}
	res = append(res, stopByte)
	return res
}

func BytesToDomain(data []byte) (int, string) {
	var res string
	var index int
	for {
		sublen := int(data[index])
		res += string(data[index+1 : index+sublen+1])
		index += (sublen + 1)
		if index < len(data) && data[index] != stopByte {
			res += "."
		} else {
			break
		}
	}
	index++
	return index, res
}

func Ipv4ToBytes(ipv4 string) ([]byte, error) {
	var res []byte
	var err error
	ipv4s := strings.Split(ipv4, ".")
	if len(ipv4s) != 4 {
		err = errors.New("Invalid ipv4 format")
		return res, err
	}
	for index := range ipv4s {
		num, err := strconv.Atoi(ipv4s[index])
		if err != nil {
			return res, err
		}
		if num > 127 {
			res = append(res, byte(num-256))
		} else {
			res = append(res, byte(num))
		}
	}
	return res, nil
}
