package option

type Option92AssociatedIpOption struct {
	TypeListUint32
}

func (_ Option92AssociatedIpOption) GetNum() byte {
	return byte(92)
}

func (this *Option92AssociatedIpOption) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option92AssociatedIpOption) Construct(associatedIpOption []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), associatedIpOption)
}

func (this Option92AssociatedIpOption) GetAssociatedIpOption() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option92AssociatedIpOption) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
