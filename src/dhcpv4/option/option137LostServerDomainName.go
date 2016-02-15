package option

type Option137LostServerDomainName struct {
	TypeString
}

func (_ Option137LostServerDomainName) GetNum() byte {
	return byte(137)
}

func (this *Option137LostServerDomainName) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option137LostServerDomainName) GetLostServerDomainName() string {
	return this.TypeString.GetString()
}

func (this *Option137LostServerDomainName) Construct(lostServerDomainName string) {
	this.TypeString.Construct(this.GetNum(), lostServerDomainName)
}

func (this Option137LostServerDomainName) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
