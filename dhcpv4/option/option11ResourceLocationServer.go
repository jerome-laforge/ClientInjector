package option

type Option11ResourceLocationServer struct {
	TypeListUint32
}

func (_ Option11ResourceLocationServer) GetNum() byte {
	return byte(11)
}

func (this *Option11ResourceLocationServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option11ResourceLocationServer) Construct(resourceLocationServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), resourceLocationServer)
}

func (this Option11ResourceLocationServer) GetResourceLocationServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option11ResourceLocationServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
