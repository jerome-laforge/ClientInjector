package option

type Option28BroadcastAddressOption struct {
	TypeUint32
}

func (_ Option28BroadcastAddressOption) GetNum() byte {
	return byte(28)
}

func (this *Option28BroadcastAddressOption) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option28BroadcastAddressOption) Construct(leaseTime uint32) {
	this.TypeUint32.Construct(this.GetNum(), leaseTime)
}

func (this Option28BroadcastAddressOption) GetBroadcastAddress() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option28BroadcastAddressOption) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
