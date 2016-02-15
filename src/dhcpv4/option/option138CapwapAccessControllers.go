package option

type Option138CapwapAccessControllers struct {
	TypeListUint32
}

func (_ Option138CapwapAccessControllers) GetNum() byte {
	return byte(138)
}

func (this *Option138CapwapAccessControllers) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option138CapwapAccessControllers) Construct(capwapAccessControllers []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), capwapAccessControllers)
}

func (this Option138CapwapAccessControllers) GetCapwapAccessControllers() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option138CapwapAccessControllers) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
