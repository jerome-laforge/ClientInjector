package option

type Option66TftpServerName struct {
	TypeString
}

func (_ Option66TftpServerName) GetNum() byte {
	return byte(66)
}

func (this *Option66TftpServerName) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option66TftpServerName) GetTftpServerName() string {
	return this.TypeString.GetString()
}

func (this *Option66TftpServerName) Construct(tftpServerName string) {
	this.TypeString.Construct(this.GetNum(), tftpServerName)
}

func (this Option66TftpServerName) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
