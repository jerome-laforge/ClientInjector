// +build gofuzz

package dhcpv4

import "github.com/jerome-laforge/ClientInjector/dhcpv4/option"

var opt60 *option.Option60VendorClassIdentifier

func init() {
	opt60 = new(option.Option60VendorClassIdentifier)
	opt60.Construct("toto")
}

//func Fuzz(roData []byte) int {
//	data := make([]byte, len(roData))
//	copy(data, roData)
func Fuzz(data []byte) int {
	errno, p, o := fuzzParse(data)
	if errno != 1 {
		return errno
	}

	if o == nil {
		o = opt60
	}
	p.AddOption(o)

	newErrno, _, _ := fuzzParse(p.Raw)
	if newErrno != errno {
		panic("there is inconsistency after adding option or reparsing the packet")
	}

	p.RemoveAllOptionsExcept([]byte{53, 82})
	p.AddOption(o)

	newErrno, p, _ = fuzzParse(p.Raw)
	if newErrno != errno {
		panic("there is inconsistency after RemoveAllOptionsExcept")
	}

	return errno
}

func fuzzParse(data []byte) (int, *DhcpPacket, option.SpecificOption) {
	var o option.SpecificOption
	p, err := Parse(data)
	if err != nil {
		if p != nil {
			panic("packet != nil on error")
		}
		return 0, nil, o
	}

	if p == nil {
		panic("packet == nil when no error")
	}

	if !p.ContainsMagicCookie() {
		panic("No magic Cookie...")
	}

	if p.GetMacAddr() == nil {
		panic("MacAddr is nil")
	}

	if len(p.GetMacAddr()) > 16 {
		panic("MacAddr is greater than 16")
	}

	{ // option 53
		opt := new(option.Option53MessageType)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				if opt.GetLength() != 1 {
					panic("len opt53 != 1")
				}
				if len(opt.GetValue()) != int(opt.GetLength()) {
					panic("len opt53 inconsistency")
				}

				o = opt
			}
		}
	}

	{ // option 51 and all TypeUint32 option type (e.g. option 54, option 58, option 59 ...)
		opt := new(option.Option51IpAddressLeaseTime)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				if opt.GetLength() != 4 {
					panic("len opt51 != 4")
				}
				if len(opt.GetValue()) != int(opt.GetLength()) {
					panic("len opt51 inconsistency")
				}

				o = opt
			}
		}
	}

	{ // option 82
		opt := new(option.Option82DhcpAgentOption)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				if len(opt.GetValue()) != int(opt.GetLength()) {
					panic("len opt82 inconsistency")
				}

				{ // option 82.1
					ssOpt := new(option.Option82_1CircuitId)
					opt.GetSubOptions(ssOpt)
				}

				{ // option 82.2
					ssOpt := new(option.Option82_2RemoteId)
					opt.GetSubOptions(ssOpt)
				}

				{ // option 82.9
					ssOpt := new(option.Option82_9VendorSpecificInformation)
					_, _ = opt.GetSubOptions(ssOpt)
				}

				o = opt
			}
		}
	}

	{ // option 60 and all TypeString option type (option 15 ...)
		opt := new(option.Option60VendorClassIdentifier)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				if len(opt.GetString()) != int(opt.GetLength()) {
					panic("len opt60 inconsistency")
				}
			}
		}
	}

	{ // option 6 and all TypeListUint32 option type (option 10, option 11 ...)
		opt := new(option.Option6DomainNameServerOption)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				if len(opt.GetValue()) != int(opt.GetLength()) {
					panic("len opt6 inconsistency")
				}
			}
		}
	}

	{ // option 119
		opt := new(option.Option119DomainSearch)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				_ = opt.GetSearchString()
			}
		}
	}

	{ // option 120
		opt := new(option.Option120SipServers)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				switch opt.GetEnc() {
				case option.ENC_FQDN:
					fqdns, e := opt.GetFqdnList()
					if e != nil && fqdns != nil {
						panic("opt120 ENC_FQDN fqdns != nil on error")
					}
				case option.ENC_IP:
					addrs, e := opt.GetIpAddresses()
					if e != nil && addrs != nil {
						panic("opt120 ENC_IP addrs != nil on error")
					}
				default:
					panic("opt120 with ENC_UNKNOWN is not possible to have error == nil on p.GetOption(opt)")
				}
			}
		}
	}

	{ // option 119
		opt := new(option.Option119DomainSearch)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				_ = opt.GetSearchString()
			}
		}
	}

	{ // option 61
		opt := new(option.Option61ClientIdentifier)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				_ = opt.GetClientIdentifier()
			}
		}
	}

	{ // option 77
		opt := new(option.Option77UserClass)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				if len(opt.GetValue()) != int(opt.GetLength()) {
					panic("len opt77 inconsistency")
				}
				for _, uc := range opt.UserClassDatas {
					if len(uc.GetUcData()) != int(uc.GetUcLen()) {
						panic("opt77 UserClass len inconsistency")
					}
				}
			}
		}
	}

	{ // option 90
		opt := new(option.Option90Authentificiation)
		if f, e := p.GetOption(opt); e == nil {
			if f {
				opt.GetProtocol()
				opt.GetAlgorithm()
				opt.GetRdm()
				opt.GetReplayDetection()
				opt.GetAuthenticationInformation()
			}
		}
	}

	{ // option 121
		opt := new(option.Option121ClasslessStaticRoute)
		p.GetOption(opt)
	}

	return 1, p, o
}
