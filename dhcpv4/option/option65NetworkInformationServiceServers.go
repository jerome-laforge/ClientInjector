package option

type Option65NetworkInformationServiceServers struct {
	TypeListUint32
}

func (_ Option65NetworkInformationServiceServers) GetNum() byte {
	return byte(65)
}

func (this *Option65NetworkInformationServiceServers) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option65NetworkInformationServiceServers) Construct(networkInformationServiceServers []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), networkInformationServiceServers)
}

func (this Option65NetworkInformationServiceServers) GetNetworkInformationServiceServers() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option65NetworkInformationServiceServers) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
