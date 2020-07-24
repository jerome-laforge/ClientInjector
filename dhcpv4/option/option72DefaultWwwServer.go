package option

type Option72DefaultWwwServer struct {
	TypeListUint32
}

func (_ Option72DefaultWwwServer) GetNum() byte {
	return byte(72)
}

func (this *Option72DefaultWwwServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option72DefaultWwwServer) Construct(defaultWwwServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), defaultWwwServer)
}

func (this Option72DefaultWwwServer) GetDefaultWwwServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option72DefaultWwwServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
