package option

type Option4TimeServerOption struct {
	TypeListUint32
}

func (_ Option4TimeServerOption) GetNum() byte {
	return byte(4)
}

func (this *Option4TimeServerOption) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option4TimeServerOption) Construct(addresses []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), addresses)
}

func (this Option4TimeServerOption) GetAddresses() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option4TimeServerOption) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
