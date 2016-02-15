package option

type Option74DefaultIrcServer struct {
	TypeListUint32
}

func (_ Option74DefaultIrcServer) GetNum() byte {
	return byte(74)
}

func (this *Option74DefaultIrcServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option74DefaultIrcServer) Construct(defaultIrcServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), defaultIrcServer)
}

func (this Option74DefaultIrcServer) GetDefaultIrcServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option74DefaultIrcServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
