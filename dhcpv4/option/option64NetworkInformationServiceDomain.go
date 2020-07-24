package option

type Option64NetworkInformationServiceDomain struct {
	TypeString
}

func (_ Option64NetworkInformationServiceDomain) GetNum() byte {
	return byte(64)
}

func (this *Option64NetworkInformationServiceDomain) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option64NetworkInformationServiceDomain) GetNetworkInformationServiceDomain() string {
	return this.TypeString.GetString()
}

func (this *Option64NetworkInformationServiceDomain) Construct(networkInformationServiceDomain string) {
	this.TypeString.Construct(this.GetNum(), networkInformationServiceDomain)
}

func (this Option64NetworkInformationServiceDomain) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
