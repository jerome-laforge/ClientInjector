package option

type Option7LogServer struct {
	TypeListUint32
}

func (_ Option7LogServer) GetNum() byte {
	return byte(7)
}

func (this *Option7LogServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option7LogServer) Construct(logServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), logServer)
}

func (this Option7LogServer) GetLogServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option7LogServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
