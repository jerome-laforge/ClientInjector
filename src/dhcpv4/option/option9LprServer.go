package option

type Option9LprServer struct {
	TypeListUint32
}

func (_ Option9LprServer) GetNum() byte {
	return byte(9)
}

func (this *Option9LprServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option9LprServer) Construct(lprServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), lprServer)
}

func (this Option9LprServer) GetLprServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option9LprServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
