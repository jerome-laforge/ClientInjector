package option

type Option58RenewalTimeValue struct {
	TypeUint32
}

func (_ Option58RenewalTimeValue) GetNum() byte {
	return byte(58)
}

func (this *Option58RenewalTimeValue) Parse(rawOpt RawOption) error {
	return this.TypeUint32.Parse(this.GetNum(), rawOpt)
}

func (this *Option58RenewalTimeValue) Construct(leaseTime uint32) {
	this.TypeUint32.Construct(this.GetNum(), leaseTime)
}

func (this Option58RenewalTimeValue) GetRenewalTime() uint32 {
	return this.TypeUint32.GetUint32()
}

func (this Option58RenewalTimeValue) GetRawOption() RawOption {
	return this.TypeUint32.GetRawOption()
}
