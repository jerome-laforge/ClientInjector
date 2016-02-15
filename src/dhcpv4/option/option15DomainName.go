package option

type Option15DomainName struct {
	TypeString
}

func (_ Option15DomainName) GetNum() byte {
	return byte(15)
}

func (this *Option15DomainName) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option15DomainName) GetDomainName() string {
	return this.TypeString.GetString()
}

func (this *Option15DomainName) Construct(domainName string) {
	this.TypeString.Construct(this.GetNum(), domainName)
}

func (this Option15DomainName) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
