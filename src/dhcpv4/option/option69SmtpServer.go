package option

type Option69SmtpServer struct {
	TypeListUint32
}

func (_ Option69SmtpServer) GetNum() byte {
	return byte(69)
}

func (this *Option69SmtpServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option69SmtpServer) Construct(smtpServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), smtpServer)
}

func (this Option69SmtpServer) GetSmtpServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option69SmtpServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
