package option

type Option87NovellDirectoryServicesContext struct {
	TypeString
}

func (_ Option87NovellDirectoryServicesContext) GetNum() byte {
	return byte(87)
}

func (this *Option87NovellDirectoryServicesContext) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option87NovellDirectoryServicesContext) GetNovellDirectoryServicesContext() string {
	return this.TypeString.GetString()
}

func (this *Option87NovellDirectoryServicesContext) Construct(novellDirectoryServicesContext string) {
	this.TypeString.Construct(this.GetNum(), novellDirectoryServicesContext)
}

func (this Option87NovellDirectoryServicesContext) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
