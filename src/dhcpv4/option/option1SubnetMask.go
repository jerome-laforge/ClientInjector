package option

type Option1SubnetMask struct {
	TypeUint32
}

func (_ Option1SubnetMask) GetNum() byte {
	return byte(1)
}

func (this *Option1SubnetMask) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option1SubnetMask) Construct(mask uint32) {
	this.TypeUint32.Construct(this.GetNum(), mask)
}

func (this Option1SubnetMask) GetMask() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option1SubnetMask) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
