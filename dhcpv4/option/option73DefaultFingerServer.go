package option

type Option73DefaultFingerServer struct {
	TypeListUint32
}

func (_ Option73DefaultFingerServer) GetNum() byte {
	return byte(73)
}

func (this *Option73DefaultFingerServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option73DefaultFingerServer) Construct(defaultFingerServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), defaultFingerServer)
}

func (this Option73DefaultFingerServer) GetDefaultFingerServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option73DefaultFingerServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
