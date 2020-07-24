package option

type Option67BootfileName struct {
	TypeString
}

func (_ Option67BootfileName) GetNum() byte {
	return byte(67)
}

func (this *Option67BootfileName) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option67BootfileName) GetBootfileName() string {
	return this.TypeString.GetString()
}

func (this *Option67BootfileName) Construct(bootfileName string) {
	this.TypeString.Construct(this.GetNum(), bootfileName)
}

func (this Option67BootfileName) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
