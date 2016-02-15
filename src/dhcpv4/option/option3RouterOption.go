package option

type Option3RouterOption struct {
	TypeListUint32
}

func (_ Option3RouterOption) GetNum() byte {
	return byte(3)
}

func (this *Option3RouterOption) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option3RouterOption) Construct(addresses []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), addresses)
}

func (this Option3RouterOption) GetAddresses() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option3RouterOption) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
