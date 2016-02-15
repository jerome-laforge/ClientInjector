package option

type Option118SubnetSelectionOption struct {
	TypeListUint32
}

func (_ Option118SubnetSelectionOption) GetNum() byte {
	return byte(118)
}

func (this *Option118SubnetSelectionOption) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option118SubnetSelectionOption) Construct(subnetSelectionOption []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), subnetSelectionOption)
}

func (this Option118SubnetSelectionOption) GetSubnetSelectionOption() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option118SubnetSelectionOption) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
