package option

type Option40NetworkInformationServiceDomain struct {
	TypeString
}

func (_ Option40NetworkInformationServiceDomain) GetNum() byte {
	return byte(40)
}

func (this *Option40NetworkInformationServiceDomain) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option40NetworkInformationServiceDomain) GetNetworkInformationServiceDomain() string {
	return this.TypeString.GetString()
}

func (this *Option40NetworkInformationServiceDomain) Construct(networkInformationServiceDomain string) {
	this.TypeString.Construct(this.GetNum(), networkInformationServiceDomain)
}

func (this Option40NetworkInformationServiceDomain) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
