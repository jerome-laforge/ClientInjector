package option

type Option54DhcpServerIdentifier struct {
	TypeUint32
}

func (_ Option54DhcpServerIdentifier) GetNum() byte {
	return byte(54)
}

func (this *Option54DhcpServerIdentifier) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option54DhcpServerIdentifier) Construct(dhcpServerIdentifier uint32) {
	this.TypeUint32.Construct(this.GetNum(), dhcpServerIdentifier)
}

func (this Option54DhcpServerIdentifier) GetDhcpServerIdentifier() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option54DhcpServerIdentifier) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
