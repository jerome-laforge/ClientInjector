package option

type Option10ImpressServer struct {
	TypeListUint32
}

func (_ Option10ImpressServer) GetNum() byte {
	return byte(10)
}

func (this *Option10ImpressServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option10ImpressServer) Construct(impressServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), impressServer)
}

func (this Option10ImpressServer) GetImpressServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option10ImpressServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
