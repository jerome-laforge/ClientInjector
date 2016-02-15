package option

type Option62NovellNetwareIpdomain struct {
	TypeString
}

func (_ Option62NovellNetwareIpdomain) GetNum() byte {
	return byte(62)
}

func (this *Option62NovellNetwareIpdomain) Parse(rawOpt RawOption) error {
	return this.TypeString.Parse(this.GetNum(), rawOpt)
}

func (this Option62NovellNetwareIpdomain) GetNovellNetwareIpdomain() string {
	return this.TypeString.GetString()
}

func (this *Option62NovellNetwareIpdomain) Construct(novellNetwareIpdomain string) {
	this.TypeString.Construct(this.GetNum(), novellNetwareIpdomain)
}

func (this Option62NovellNetwareIpdomain) GetRawOption() RawOption {
	return this.TypeString.GetRawOption()
}
