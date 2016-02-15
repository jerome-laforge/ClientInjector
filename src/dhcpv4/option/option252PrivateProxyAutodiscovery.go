package option

type Option252PrivateProxyAutodiscovery struct {
	TypeString
}

func (_ Option252PrivateProxyAutodiscovery) GetNum() byte {
	return byte(252)
}

func (this *Option252PrivateProxyAutodiscovery) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option252PrivateProxyAutodiscovery) GetProxyAutodiscovery() string {
	return this.TypeString.GetString()
}

func (this *Option252PrivateProxyAutodiscovery) Construct(proxyAutodiscovery string) {
	this.TypeString.Construct(this.GetNum(), proxyAutodiscovery)
}

func (this Option252PrivateProxyAutodiscovery) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
