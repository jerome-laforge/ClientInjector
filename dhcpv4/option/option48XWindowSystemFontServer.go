package option

type Option48XWindowSystemFontServer struct {
	TypeListUint32
}

func (_ Option48XWindowSystemFontServer) GetNum() byte {
	return byte(48)
}

func (this *Option48XWindowSystemFontServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option48XWindowSystemFontServer) Construct(xWindowSystemFontServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), xWindowSystemFontServer)
}

func (this Option48XWindowSystemFontServer) GetXWindowSystemFontServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option48XWindowSystemFontServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
