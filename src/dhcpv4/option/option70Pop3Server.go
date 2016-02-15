package option

type Option70Pop3Server struct {
	TypeListUint32
}

func (_ Option70Pop3Server) GetNum() byte {
	return byte(70)
}

func (this *Option70Pop3Server) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option70Pop3Server) Construct(pop3Server []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), pop3Server)
}

func (this Option70Pop3Server) GetPop3Server() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option70Pop3Server) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
