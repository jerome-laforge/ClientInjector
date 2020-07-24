package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jerome-laforge/ClientInjector/dhcpv4/dherrors"
)

func Convert4byteToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func ConvertUint32ToNew4byte(num uint32) []byte {
	ipRaw := make([]byte, 4)
	ConvertUint32To4byte(num, ipRaw)
	return ipRaw
}

func ConvertUint32To4byte(num uint32, b []byte) {
	binary.BigEndian.PutUint32(b, num)
}

func Convert8byteToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func ConvertMax8byteToUint64(b []byte) (ret uint64) {
	l := len(b)
	if l > 8 {
		l = 8
	}

	for i := 0; i < l; i++ {
		ret = ret << 8
		ret = ret | uint64(b[i])
	}
	return
}

func ConvertUint64To8byte(num uint64, b []byte) {
	binary.BigEndian.PutUint64(b, num)
}

func ConvertUint32ToIpAddr(num uint32) string {
	return strconv.Itoa(int((num&0xff000000)>>24)) + "." +
		strconv.Itoa(int((num&0x00ff0000)>>16)) + "." +
		strconv.Itoa(int((num&0x0000ff00)>>8)) + "." +
		strconv.Itoa(int(num&0x000000ff))
}

func ConvertIpAddrToUint32(ipaddr string) (uint32, error) {
	bytes := strings.Split(ipaddr, ".")
	if len(bytes) != 4 {
		return 0, errors.New(fmt.Sprintf("Bad format IPv4 address : %s", ipaddr))
	}

	var ipNum uint32 = 0

	b, err := strconv.Atoi(bytes[0])
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Bad format IPv4 address : %s", ipaddr))
	}
	ipNum += uint32(b) << 24

	b, err = strconv.Atoi(bytes[1])
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Bad format IPv4 address : %s", ipaddr))
	}
	ipNum += uint32(b) << 16

	b, err = strconv.Atoi(bytes[2])
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Bad format IPv4 address : %s", ipaddr))
	}
	ipNum += uint32(b) << 8

	b, err = strconv.Atoi(bytes[3])
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Bad format IPv4 address : %s", ipaddr))
	}
	ipNum += uint32(b)

	return ipNum, nil
}

func ConvertByteToString(data []byte) string {
	s, _ := strconv.Unquote(fmt.Sprintf("%+q", data))
	return s
}

func ConvertByteToHexString(data []byte) string {
	return fmt.Sprintf("%x", data)
}

func GetDnsName(roBuf []byte) ([]string, error) {
	// For avoiding to corrupt rawData
	// We copy option's buffer into new buffer where we can replace Null char (0x00) by . (0x2e)
	// and/or managed the compression pointer (0xc0)
	rwBuf := make([]byte, len(roBuf))
	copy(rwBuf, roBuf)
	rawFqdns := bytes.Split(rwBuf, []byte{0})
	if len(rawFqdns[len(rawFqdns)-1]) == 0 {
		rawFqdns = rawFqdns[:len(rawFqdns)-1]
	}
	fqdns := make([]string, len(rawFqdns))
	for i := len(rawFqdns) - 1; i >= 0; i-- {
		gap := 0
		j := 0
		compressFound := false
		for ; j < len(rawFqdns[i]); j += gap {
			if j >= len(rawFqdns[i]) {
				return nil, dherrors.MalformedOption
			}
			if rawFqdns[i][j] == 0xc0 { // compression pointer
				if compressFound {
					// Only once compression pointer is allowed.
					// This avoid endless loop and so dead lock with memory leak...
					return nil, dherrors.MalformedOption
				}
				compressFound = true
				if len(rawFqdns[i]) == j+1 {
					return nil, dherrors.MalformedOption
				}
				if int(rawFqdns[i][j+1]) >= len(rawFqdns[0]) {
					return nil, dherrors.MalformedOption
				}
				rawFqdns[i] = append(rawFqdns[i][:j], rawFqdns[0][int(rawFqdns[i][j+1]):]...)
			}
			gap = int(rawFqdns[i][j]) + 1
			rawFqdns[i][j] = '.'
		}
		if j != len(rawFqdns[i]) {
			return nil, errors.New(fmt.Sprintf("option is malformed [expected: %d] [actual: %d]", len(rawFqdns[i]), j))
		}
		if len(rawFqdns[i]) == 0 {
			return nil, dherrors.MalformedOption
		}
		fqdns[i] = string(rawFqdns[i][1:])
	}

	return fqdns, nil
}

func SerializeDnsName(fqdns []string) []byte {
	i := 0
	for _, fqdn := range fqdns {
		i += len(fqdn) + 2
	}
	data := make([]byte, i)
	i = 0
	for _, fqdn := range fqdns {
		slice := data[i : i+len(fqdn)+2]
		i += len(fqdn) + 2
		slice[0] = '.'
		copy(slice[1:], fqdn)
		cur := len(slice) - 1
		prev := cur
		for ; cur >= 0; cur-- {
			if slice[cur] == '.' {
				slice[cur] = byte(prev - cur - 1)
				prev = cur
			}
		}
	}
	return data
}
