package option

import (
	"errors"
	"fmt"

	"github.com/jerome-laforge/ClientInjector/dhcpv4/dherrors"
	"github.com/jerome-laforge/ClientInjector/dhcpv4/util"
)

const HEADER_LEN_OPT_120 = 3

const (
	ENC_OPT_120_FQDN    = iota
	ENC_OPT_120_IP_LIST = iota
	ENC_OPT_120_UNKNOWN = 255
)

type EncOption120 struct {
	EncOption120 string
	Value        byte
}

var ENC_FQDN EncOption120 = EncOption120{
	EncOption120: "FQDN",
	Value:        ENC_OPT_120_FQDN,
}

var ENC_IP EncOption120 = EncOption120{
	EncOption120: "IPv4 Address List",
	Value:        ENC_OPT_120_IP_LIST,
}

var ENC_UNKNOWN EncOption120 = EncOption120{
	EncOption120: "Unknow encoding type for option 120",
	Value:        ENC_OPT_120_IP_LIST,
}

type Option120SipServers struct {
	rawOpt RawOption
}

func (_ Option120SipServers) GetNum() byte {
	return byte(120)
}

func (this *Option120SipServers) Parse(rawOpt RawOption) error {
	if rawOpt.GetLength() < HEADER_LEN_OPT_120 {
		return dherrors.MalformedOption
	}
	this.rawOpt = rawOpt
	enc := this.GetEnc()
	if enc == ENC_UNKNOWN {
		return errors.New(fmt.Sprintf("Unknow encoding type for option 120 [actual: %s]", enc.EncOption120))
	}
	if enc == ENC_IP && this.rawOpt.GetLength()&3 != 0 {
		return dherrors.Opt120HasInvalidLen
	}
	return nil
}

func (this Option120SipServers) GetRawOption() RawOption {
	return this.rawOpt
}

func (this Option120SipServers) GetEnc() EncOption120 {
	enc := this.rawOpt.GetValue()[0]
	if enc == ENC_OPT_120_FQDN {
		return ENC_FQDN
	} else if enc == ENC_OPT_120_IP_LIST {
		return ENC_IP
	} else {
		return ENC_UNKNOWN
	}
}

func (this Option120SipServers) GetFqdnList() ([]string, error) {
	if this.GetEnc() != ENC_FQDN {
		return nil, errors.New(fmt.Sprintf("Enc for option 120 is not expected type [expected: %s] [actual: %s]", ENC_FQDN.EncOption120, this.GetEnc().EncOption120))
	}
	return util.GetDnsName(this.rawOpt.GetValue()[1:])
}

func (this Option120SipServers) GetIpAddresses() ([]uint32, error) {
	if this.GetEnc() != ENC_IP {
		return nil, errors.New(fmt.Sprintf("Enc for option 120 is not expected type [expected: %s] [actual: %s]", ENC_IP.EncOption120, this.GetEnc().EncOption120))
	}
	ipAddresses := make([]uint32, (this.rawOpt.GetLength()-1)>>2)
	for i := range ipAddresses {
		off := 1 + i*4
		ipAddresses[i] = util.Convert4byteToUint32(this.rawOpt.GetValue()[off : off+4])
	}
	return ipAddresses, nil
}

func (this *Option120SipServers) construct(enc EncOption120, data []byte) error {
	if enc != ENC_FQDN && enc != ENC_IP {
		return dherrors.Opt120UnsupportedEncoding
	}
	if len(data) > 255 || (enc == ENC_IP && len(data)&3 != 0) {
		return dherrors.Opt120DataInvalidLen
	}
	this.rawOpt.Construct(this.GetNum(), byte(len(data)+1))
	this.rawOpt.GetValue()[0] = enc.Value
	copy(this.rawOpt.GetValue()[1:], data)
	return nil
}

func (this *Option120SipServers) ConstructFqdn(fqdns []string) error {
	data := util.SerializeDnsName(fqdns)
	return this.construct(ENC_FQDN, data)
}

func (this *Option120SipServers) ConstructIp(ipList []uint32) error {
	data := make([]byte, 4*len(ipList))
	for i, ip := range ipList {
		util.ConvertUint32To4byte(ip, data[i*4:(i+1)*4])
	}
	return this.construct(ENC_IP, data)
}
