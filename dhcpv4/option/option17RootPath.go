package option

type Option17RootPath struct {
	TypeString
}

func (_ Option17RootPath) GetNum() byte {
	return byte(17)
}

func (this *Option17RootPath) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option17RootPath) GetRootPath() string {
	return this.TypeString.GetString()
}

func (this *Option17RootPath) Construct(rootPath string) {
	this.TypeString.Construct(this.GetNum(), rootPath)
}

func (this Option17RootPath) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
