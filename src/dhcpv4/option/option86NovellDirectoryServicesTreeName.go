package option

type Option86NovellDirectoryServicesTreeName struct {
	TypeString
}

func (_ Option86NovellDirectoryServicesTreeName) GetNum() byte {
	return byte(86)
}

func (this *Option86NovellDirectoryServicesTreeName) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option86NovellDirectoryServicesTreeName) GetNovellDirectoryServicesTreeName() string {
	return this.TypeString.GetString()
}

func (this *Option86NovellDirectoryServicesTreeName) Construct(novellDirectoryServicesTreeName string) {
	this.TypeString.Construct(this.GetNum(), novellDirectoryServicesTreeName)
}

func (this Option86NovellDirectoryServicesTreeName) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
