package option

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jerome-laforge/ClientInjector/dhcpv4/dherrors"
	"github.com/jerome-laforge/ClientInjector/dhcpv4/util"
)

const (
	HEADER_LEN_OPT_121 = 5
	MAX_CIDR           = 32
)

type Option121ClasslessStaticRoute struct {
	rawOpt RawOption
	Routes []Route
}

type Route struct {
	Dest uint32
	Mask byte
	Gw   uint32
}

func (this Route) String() string {
	return "Subnet/MaskWidth : " + util.ConvertUint32ToIpAddr(this.Dest) + "/" + strconv.Itoa(int(this.Mask)) + " Router : " + util.ConvertUint32ToIpAddr(this.Gw)
}

func (_ Option121ClasslessStaticRoute) GetNum() byte {
	return byte(121)
}

func (this *Option121ClasslessStaticRoute) Parse(rawOpt RawOption) error {
	this.rawOpt = rawOpt
	length := int(this.rawOpt.GetLength())
	if length < HEADER_LEN_OPT_121 {
		return errors.New(fmt.Sprintf("Option 121 is malformed : has to be greater than 5 bytes [actual: %d]", length))
	}

	idx := 0
	buf := make([]byte, 4)
	for idx < length {
		if this.rawOpt.GetValue()[idx] > MAX_CIDR {
			return errors.New(fmt.Sprintf("Option 121 malformed for route #%d, [expected CIDR mask lower than : %d] [actual : %d]", len(this.Routes), MAX_CIDR+1, this.rawOpt.GetValue()[idx]))
		}
		significantOctets := significantBytes(this.rawOpt.GetValue()[idx])
		if idx+significantOctets+5 > length {
			return errors.New(fmt.Sprintf("Option 121 not enough room for route #%d, [expected : %d] [actual : %d]", len(this.Routes), significantOctets+5, length-idx))
		}
		curRoute := new(Route)

		// get cidr
		curRoute.Mask = this.rawOpt.GetValue()[idx]

		// get dest
		copy(buf, this.rawOpt.GetValue()[idx+1:idx+significantOctets+1])
		curRoute.Dest = util.Convert4byteToUint32(buf)

		// get gw
		curRoute.Gw = util.Convert4byteToUint32(this.rawOpt.GetValue()[idx+significantOctets+1 : idx+significantOctets+5])
		this.Routes = append(this.Routes, *curRoute)

		idx += significantOctets + 5
		for i := range buf {
			buf[i] = 0
		}
	}
	return nil
}

func (this *Option121ClasslessStaticRoute) Construct(routes []Route) error {
	buf := make([]byte, 256)
	idx := 0
	for _, route := range routes {
		if 5+significantBytes(route.Mask)+idx >= len(buf) {
			return dherrors.Opt121TooLong
		}
		buf[idx] = route.Mask
		util.ConvertUint32To4byte(route.Dest, buf[idx+1:idx+5])
		util.ConvertUint32To4byte(route.Gw, buf[idx+1+significantBytes(route.Mask):idx+5+significantBytes(route.Mask)])
		idx += 5 + significantBytes(route.Mask)
	}

	buf = buf[:idx]
	this.rawOpt.Construct(this.GetNum(), byte(len(buf)))
	copy(this.rawOpt.GetValue(), buf)
	return nil
}

func (this Option121ClasslessStaticRoute) GetRawOption() RawOption {
	return this.rawOpt
}

func significantBytes(mask byte) int {
	if mask&7 == 0 {
		return int(mask >> 3)
	}
	return int(mask>>3 + 1)
}
