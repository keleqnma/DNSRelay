package common

const stopByte = 0x00

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
	var res []byte
	res = append(res, byte((m>>8)&0xFF))
	res = append(res, byte((m)&0xFF))
	return res
}

func IntToBytes(num int) []byte {
	m := int32(num)
	var res []byte
	res = append(res, byte((m>>24)&0xFF))
	res = append(res, byte((m>>16)&0xFF))
	res = append(res, byte((m>>8)&0xFF))
	res = append(res, byte((m)&0xFF))
	return res
}

func BytesToInt8(b []byte) int8 {
	return int8(((b[0] & 0xff) << 8) | (b[1] & 0xff))
}

func DomainToBytes(domain string) []byte {
	var res []byte
	// nums := strings.Split(domain, ".")
	// for index := range nums {
	// 	res = append(res, []byte(nums[index])...)
	// }
	res = append(res, []byte(domain)...)
	res = append(res, stopByte)
	return res
}

func BytesToDomain(data []byte) (int, string) {
	var res string
	var index int
	for index = 0; index < len(data) && data[index] != stopByte; index++ {
		res += string(data[index])
	}
	index++
	return index, res
}
