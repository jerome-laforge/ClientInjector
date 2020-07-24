package option

type Option41NetworkInformationServiceServers struct {
	TypeListUint32
}

func (_ Option41NetworkInformationServiceServers) GetNum() byte {
	return byte(41)
}

func (this *Option41NetworkInformationServiceServers) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option41NetworkInformationServiceServers) Construct(networkInformationServiceServers []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), networkInformationServiceServers)
}

func (this Option41NetworkInformationServiceServers) GetNetworkInformationServiceServers() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option41NetworkInformationServiceServers) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
