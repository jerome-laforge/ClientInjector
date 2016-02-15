package option

type Option16SwapServer struct {
	TypeUint32
}

func (_ Option16SwapServer) GetNum() byte {
	return byte(16)
}

func (this *Option16SwapServer) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option16SwapServer) Construct(swapServer uint32) {
	this.TypeUint32.Construct(this.GetNum(), swapServer)
}

func (this Option16SwapServer) GetSwapServer() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option16SwapServer) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
