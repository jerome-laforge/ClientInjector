package option

type Option47NetBiosOverTcpipScope struct {
	TypeString
}

func (_ Option47NetBiosOverTcpipScope) GetNum() byte {
	return byte(47)
}

func (this *Option47NetBiosOverTcpipScope) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option47NetBiosOverTcpipScope) GetNetBiosOverTcpipScope() string {
	return this.TypeString.GetString()
}

func (this *Option47NetBiosOverTcpipScope) Construct(netBiosOverTcpipScope string) {
	this.TypeString.Construct(this.GetNum(), netBiosOverTcpipScope)
}

func (this Option47NetBiosOverTcpipScope) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
