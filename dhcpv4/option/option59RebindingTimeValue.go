package option

type Option59RebindingTimeValue struct {
	TypeUint32
}

func (_ Option59RebindingTimeValue) GetNum() byte {
	return byte(59)
}

func (this *Option59RebindingTimeValue) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option59RebindingTimeValue) Construct(leaseTime uint32) {
	this.TypeUint32.Construct(this.GetNum(), leaseTime)
}

func (this Option59RebindingTimeValue) GetRebindingTime() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option59RebindingTimeValue) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
