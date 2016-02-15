package option

type Option42NetworkTimeProtocolServersOption struct {
	TypeListUint32
}

func (_ Option42NetworkTimeProtocolServersOption) GetNum() byte {
	return byte(42)
}

func (this *Option42NetworkTimeProtocolServersOption) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option42NetworkTimeProtocolServersOption) Construct(addresses []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), addresses)
}

func (this Option42NetworkTimeProtocolServersOption) GetAddresses() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option42NetworkTimeProtocolServersOption) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
