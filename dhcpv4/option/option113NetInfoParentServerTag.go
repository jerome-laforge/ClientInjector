package option

type Option113NetInfoParentServerTag struct {
	TypeString
}

func (_ Option113NetInfoParentServerTag) GetNum() byte {
	return byte(113)
}

func (this *Option113NetInfoParentServerTag) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option113NetInfoParentServerTag) GetNetInfoParentServerTag() string {
	return this.TypeString.GetString()
}

func (this *Option113NetInfoParentServerTag) Construct(netInfoParentServerTag string) {
	this.TypeString.Construct(this.GetNum(), netInfoParentServerTag)
}

func (this Option113NetInfoParentServerTag) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
