package option

type Option51IpAddressLeaseTime struct {
	TypeUint32
}

func (_ Option51IpAddressLeaseTime) GetNum() byte {
	return byte(51)
}

func (this *Option51IpAddressLeaseTime) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option51IpAddressLeaseTime) Construct(leaseTime uint32) {
	this.TypeUint32.Construct(this.GetNum(), leaseTime)
}

func (this Option51IpAddressLeaseTime) GetLeaseTime() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option51IpAddressLeaseTime) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
