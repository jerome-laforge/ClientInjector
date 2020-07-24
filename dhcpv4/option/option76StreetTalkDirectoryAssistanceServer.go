package option

type Option76StreetTalkDirectoryAssistanceServer struct {
	TypeListUint32
}

func (_ Option76StreetTalkDirectoryAssistanceServer) GetNum() byte {
	return byte(76)
}

func (this *Option76StreetTalkDirectoryAssistanceServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option76StreetTalkDirectoryAssistanceServer) Construct(streetTalkDirectoryAssistanceServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), streetTalkDirectoryAssistanceServer)
}

func (this Option76StreetTalkDirectoryAssistanceServer) GetStreetTalkDirectoryAssistanceServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option76StreetTalkDirectoryAssistanceServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
