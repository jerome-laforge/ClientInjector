package option

type Option71NntpServer struct {
	TypeListUint32
}

func (_ Option71NntpServer) GetNum() byte {
	return byte(71)
}

func (this *Option71NntpServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option71NntpServer) Construct(nntpServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), nntpServer)
}

func (this Option71NntpServer) GetNntpServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option71NntpServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
