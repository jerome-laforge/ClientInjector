package option

type Option50RequestedIpAddress struct {
	TypeUint32
}

func (_ Option50RequestedIpAddress) GetNum() byte {
	return byte(50)
}

func (this *Option50RequestedIpAddress) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option50RequestedIpAddress) Construct(requestedIpAddress uint32) {
	this.TypeUint32.Construct(this.GetNum(), requestedIpAddress)
}

func (this Option50RequestedIpAddress) GetRequestedIpAddress() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option50RequestedIpAddress) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
