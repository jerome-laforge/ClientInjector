package option

type Option14MeritDumpFile struct {
	TypeString
}

func (_ Option14MeritDumpFile) GetNum() byte {
	return byte(14)
}

func (this *Option14MeritDumpFile) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option14MeritDumpFile) GetMeritDumpFile() string {
	return this.TypeString.GetString()
}

func (this *Option14MeritDumpFile) Construct(meritDumpFile string) {
	this.TypeString.Construct(this.GetNum(), meritDumpFile)
}

func (this Option14MeritDumpFile) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
