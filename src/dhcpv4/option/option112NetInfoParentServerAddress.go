package option

type Option112NetInfoParentServerAddress struct {
	TypeListUint32
}

func (_ Option112NetInfoParentServerAddress) GetNum() byte {
	return byte(112)
}

func (this *Option112NetInfoParentServerAddress) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option112NetInfoParentServerAddress) Construct(netInfoParentServerAddress []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), netInfoParentServerAddress)
}

func (this Option112NetInfoParentServerAddress) GetNetInfoParentServerAddress() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option112NetInfoParentServerAddress) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
