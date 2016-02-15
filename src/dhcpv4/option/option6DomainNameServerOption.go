package option

type Option6DomainNameServerOption struct {
	TypeListUint32
}

func (_ Option6DomainNameServerOption) GetNum() byte {
	return byte(6)
}

func (this *Option6DomainNameServerOption) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option6DomainNameServerOption) Construct(addresses []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), addresses)
}

func (this Option6DomainNameServerOption) GetAddresses() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option6DomainNameServerOption) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
