package option

type Option32RouterSolicitationAddress struct {
	TypeUint32
}

func (_ Option32RouterSolicitationAddress) GetNum() byte {
	return byte(32)
}

func (this *Option32RouterSolicitationAddress) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option32RouterSolicitationAddress) Construct(routerSolicitationAddress uint32) {
	this.TypeUint32.Construct(this.GetNum(), routerSolicitationAddress)
}

func (this Option32RouterSolicitationAddress) GetRouterSolicitationAddress() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option32RouterSolicitationAddress) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
