package option

type Option12HostName struct {
	TypeString
}

func (_ Option12HostName) GetNum() byte {
	return byte(12)
}

func (this *Option12HostName) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option12HostName) GetHostName() string {
	return this.TypeString.GetString()
}

func (this *Option12HostName) Construct(hostName string) {
	this.TypeString.Construct(this.GetNum(), hostName)
}

func (this Option12HostName) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
