package option

type Option18ExtensionsPath struct {
	TypeString
}

func (_ Option18ExtensionsPath) GetNum() byte {
	return byte(18)
}

func (this *Option18ExtensionsPath) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option18ExtensionsPath) GetExtensionsPath() string {
	return this.TypeString.GetString()
}

func (this *Option18ExtensionsPath) Construct(extensionsPath string) {
	this.TypeString.Construct(this.GetNum(), extensionsPath)
}

func (this Option18ExtensionsPath) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
