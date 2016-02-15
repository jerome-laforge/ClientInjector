package option

type Option45NetBiosTcpipDatagramDistributionNameServer struct {
	TypeListUint32
}

func (_ Option45NetBiosTcpipDatagramDistributionNameServer) GetNum() byte {
	return byte(45)
}

func (this *Option45NetBiosTcpipDatagramDistributionNameServer) Parse(rawOpt RawOption) error {
	return this.TypeListUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option45NetBiosTcpipDatagramDistributionNameServer) Construct(netBiosTcpipDatagramDistributionNameServer []uint32) {
	this.TypeListUint32.Construct(this.GetNum(), netBiosTcpipDatagramDistributionNameServer)
}

func (this Option45NetBiosTcpipDatagramDistributionNameServer) GetNetBiosTcpipDatagramDistributionNameServer() []uint32 {
	return this.TypeListUint32.GetListUint32()
}

func (this Option45NetBiosTcpipDatagramDistributionNameServer) GetRawOption() RawOption {
	return this.TypeListUint32.GetRawOption()
}
