package option

type Option150TftpServerAddress struct {
	TypeListUint32
}

func (_ Option150TftpServerAddress) GetNum() byte {
	return byte(150)
}

func (this *Option150TftpServerAddress) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option150TftpServerAddress) Construct(tftpServerAddress []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), tftpServerAddress)
}

func (this Option150TftpServerAddress) GetTftpServerAddress() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option150TftpServerAddress) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
