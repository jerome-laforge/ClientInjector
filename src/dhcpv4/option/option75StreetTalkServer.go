package option

type Option75StreetTalkServer struct {
	TypeListUint32
}

func (_ Option75StreetTalkServer) GetNum() byte {
	return byte(75)
}

func (this *Option75StreetTalkServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option75StreetTalkServer) Construct(streetTalkServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), streetTalkServer)
}

func (this Option75StreetTalkServer) GetStreetTalkServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option75StreetTalkServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
