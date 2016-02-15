package option

type Option5NameServer struct {
	TypeListUint32
}

func (_ Option5NameServer) GetNum() byte {
	return byte(5)
}

func (this *Option5NameServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option5NameServer) Construct(nameServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), nameServer)
}

func (this Option5NameServer) GetNameServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option5NameServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
