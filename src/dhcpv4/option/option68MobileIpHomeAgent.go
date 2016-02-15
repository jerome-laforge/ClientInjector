package option

type Option68MobileIpHomeAgent struct {
	TypeListUint32
}

func (_ Option68MobileIpHomeAgent) GetNum() byte {
	return byte(68)
}

func (this *Option68MobileIpHomeAgent) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option68MobileIpHomeAgent) Construct(mobileIpHomeAgent []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), mobileIpHomeAgent)
}

func (this Option68MobileIpHomeAgent) GetMobileIpHomeAgent() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option68MobileIpHomeAgent) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
